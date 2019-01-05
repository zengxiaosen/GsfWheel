package worker

import (
	"context"
	"go.etcd.io/etcd/clientv3"
	"gsf/src/DistributeCrotab/common"
)

//分布式锁(TXN事务)
type JobLock struct {
	//etcd客户端
	kv clientv3.KV
	//租约api，通过lease实现自动过期，避免节点宕机后锁永远占用
	lease clientv3.Lease
	//任务名
	jobName string
	//用于终止自动续租
	cancelFunc context.CancelFunc
	//租约ID
	leaseId clientv3.LeaseID
	//是否上锁成功
	isLocked bool
}

//初始化一把锁
func InitJobLock(jobName string, kv clientv3.KV, lease clientv3.Lease) (jobLock *JobLock) {
	jobLock = &JobLock{
		kv:      kv,
		lease:   lease,
		jobName: jobName,
	}
	return
}

//尝试上锁
func (jobLock *JobLock) TryLock() (err error) {

	var (
		leaseGrantResp *clientv3.LeaseGrantResponse
		cancelCtx      context.Context
		cancelFunc     context.CancelFunc
		leaseId        clientv3.LeaseID
		keepRespChan   <-chan *clientv3.LeaseKeepAliveResponse
		txn            clientv3.Txn
		lockKey        string
		txnResp        *clientv3.TxnResponse
	)
	//1,创建租约(5秒) 节点宕机，锁自动释放
	if leaseGrantResp, err = jobLock.lease.Grant(context.TODO(), 5); err != nil {
		return
	}

	//context用于取消自动续租
	cancelCtx, cancelFunc = context.WithCancel(context.TODO())

	//租约ID
	leaseId = leaseGrantResp.ID

	//2,自动续租
	if keepRespChan, err = jobLock.lease.KeepAlive(cancelCtx, leaseId); err != nil {
		goto FAIL
	}

	//3,处理续租应答的协程
	go func() {
		var (
			keepResp *clientv3.LeaseKeepAliveResponse
		)

		for {
			select {
			case keepResp = <-keepRespChan:
				//自动续租应答
				if keepResp == nil {
					//说明自动续租被取消掉了,就是被cancel了，可能是释放锁的时候，也可能是异常的时候
					goto END
				}

			}
		}

	END:
	}()

	//4,创建事务txn
	txn = jobLock.kv.Txn(context.TODO())

	//锁路径
	lockKey = common.JOB_LOCK_DIR + jobLock.jobName

	//5,事务枪锁
	//锁路径的创建版本＝0，则抢占,如果已经被占用，就get一下，也没什么用，就放在这里
	txn.If(clientv3.Compare(clientv3.CreateRevision(lockKey), "=", 0)).
		Then(clientv3.OpPut(lockKey, "", clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet(lockKey))

	//提交事务 提交事务失败是有可能抢到锁的，但是应答回来由于网络原因超时，这种时候我们就释放租约,把租约释放掉了，key马上就被删掉了
	if txnResp, err = txn.Commit(); err != nil {
		goto FAIL
	}

	//6,成功返回，失败释放租约
	if !txnResp.Succeeded {
		//锁被占用
		err = common.ERR_LOCK_ALREADY_REQUIRED
		goto FAIL
	}

	//抢锁成功
	jobLock.leaseId = leaseId
	jobLock.cancelFunc = cancelFunc
	jobLock.isLocked = true
	return

FAIL:
	//取消自动续租
	cancelFunc()
	//释放租约
	jobLock.lease.Revoke(context.TODO(), leaseId)

	return
}

//释放锁
func (jobLock *JobLock) Unlock() {
	if jobLock.isLocked {
		//取消程序自动续租的协程
		jobLock.cancelFunc()
		//释放租约
		jobLock.lease.Revoke(context.TODO(), jobLock.leaseId)
	}
}
