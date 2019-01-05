package main

import (
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"golang.org/x/net/context"
	"log"
	"time"
)



func main() {
	
	//etcdTest1()
	//etcdTest2()
	//etcdTest3()
	//etcdTest4()
	//etcdTest5()
	//etcdTest6()
	//etcdTest7()
	//etcdTest8()
	//etcdTest9()
	etcdTest10()


}

func etcdTest10(){
	//lease实现锁自动过期（宕机）
	//op操作
	//txn事务 if else then

	//删除
	var (
		config clientv3.Config
		client *clientv3.Client
		err error
		lease clientv3.Lease
		leaseGrantResp *clientv3.LeaseGrantResponse
		leaseId clientv3.LeaseID
		keepRespChan <-chan *clientv3.LeaseKeepAliveResponse
		keepResp *clientv3.LeaseKeepAliveResponse
		ctx context.Context
		cancelFunc context.CancelFunc
		kv clientv3.KV
		txn clientv3.Txn
		txnResp *clientv3.TxnResponse
	)

	config = clientv3.Config{
		Endpoints:[]string{"127.0.0.1:2379"},
		DialTimeout:5 * time.Second,
	}

	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}


	//1 上锁(创建租约，自动续租，拿到租约去抢占一个key，同时key会关联着租约)

	lease = clientv3.NewLease(client)
	//创建一个5秒的租约
	if leaseGrantResp, err = lease.Grant(context.TODO(), 5); err != nil {
		fmt.Println(err)
		return
	}
	//拿到租约ID
	leaseId = leaseGrantResp.ID

	//准备一个用于取消自动续租的context,它的父context是todo，其实没有任何用
	ctx, cancelFunc = context.WithCancel(context.TODO())

	//确保函数退出后，自动续租的携程会终止
	defer cancelFunc()

	//revoke函数是立即把租约告诉etcd释放掉，key就被删除掉，锁就被释放掉了
	defer lease.Revoke(context.TODO(), leaseId)

	//keepalive会启动一个协程序，自动续租,它的response会定期的回一个
	if keepRespChan, err = lease.KeepAlive(ctx, leaseId); err != nil {
		fmt.Println(err)
		return
	}

	//启动一个协程，去消费keepalive里面的应答
	//处理续租应答的协程
	go func() {
		for {
			select {
			case keepResp = <- keepRespChan:
				//如果续租的过程种出现了异常，也就是sdk和etcd失联了很久
				if keepRespChan == nil {
					fmt.Println("租约已经失效了")
					//结束整个协程
					goto END
				} else{
					//每秒会续租一次，所以就会收到一次应答
					fmt.Println("续租成功了,收到自动续租应答：", keepResp.ID)
				}
			}
		}
	END:
	}()


	//拿着租约去抢占一个key，if 不存在key, then设置它，else抢锁失败。用到事务
	kv = clientv3.NewKV(client)

	//创建事务
	txn = kv.Txn(context.TODO())

	//定义事务 我们希望job9的创建版本等于0
	//如果key不存在
	//要带上list，宕机后会不续租
	txn.If(clientv3.Compare(clientv3.CreateRevision("/cron/lock/job10"), "=", 0)).
		Then(clientv3.OpPut("/cron/lock/job10", "xxx", clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet("/cron/lock/job10"))//否则抢锁失败,这里去get一下，把它的值取回来

	//提交事务
	if txnResp, err = txn.Commit(); err != nil {
		fmt.Println(err)
		return
	}

	//判断是否抢到了锁
	if !txnResp.Succeeded {
		fmt.Println("锁被占用:", string(txnResp.Responses[0].GetResponseRange().Kvs[0].Value))
		return
	}

	//2 处理业务

	fmt.Println("处理任务")
	time.Sleep(100 * time.Second)


	//3 释放锁（取消自动续租，释放租约）
	//defer会把租约给释放掉,关联的kv就被删除了








}


func etcdTest9() {

	//使用op方法
	var (
		config clientv3.Config
		client *clientv3.Client
		err error
		kv clientv3.KV
		putOp clientv3.Op
		getOp clientv3.Op
		opResp clientv3.OpResponse
	)

	config = clientv3.Config{
		Endpoints:[]string{"127.0.0.1:2379"},
		DialTimeout:5 * time.Second,
	}

	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}

	kv = clientv3.NewKV(client)

	//etcd把所有方法进行抽象Op:operation，比如通过OpPut方法可以生成put操作，返回的是一个Op对象，表示一个Op操作
	//创建Op:opration
	putOp = clientv3.OpPut("/cron/jobs/job8", "")
	//执行OP
	if opResp, err = kv.Do(context.TODO(), putOp); err != nil {
		fmt.Println(err)
		return
	}

	//opResp.Put()可以从opResp 转化为 putResp
	fmt.Println("写入Revision:", opResp.Put().Header.Revision)

	//创建Op
	getOp = clientv3.OpGet("/cron/jobs/job8")

	//执行OP
	if opResp, err = kv.Do(context.TODO(), getOp); err != nil {
		fmt.Println(err)
		return
	}

	//打印
	fmt.Println("数据Revision:", opResp.Get().Kvs[0].ModRevision)
	fmt.Println("数据value:", string(opResp.Get().Kvs[0].Value))



}

func etcdTest8() {
	var (
		config clientv3.Config
		client *clientv3.Client
		err error
		kv clientv3.KV
		getResp *clientv3.GetResponse
		watchStartRevision int64
		watcher clientv3.Watcher
		watchRespChan <- chan clientv3.WatchResponse
		watchResp clientv3.WatchResponse
		event *clientv3.Event
	)

	config = clientv3.Config{
		Endpoints:[]string{"127.0.0.1:2379"},
		DialTimeout:5 * time.Second,
	}

	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}

	kv = clientv3.NewKV(client)

	//模拟etcd中kv的变化
	go func() {
		for{
			kv.Put(context.TODO(), "/cron/jobs/job8", "I am job8")
			kv.Delete(context.TODO(), "/cron/jobs/job8")
			time.Sleep(1 * time.Second)
		}
	}()
	
	//先监听到当前的值，并监听后续变化
	if getResp, err = kv.Get(context.TODO(), "/cron/jobs/job7"); err != nil {
		fmt.Println(err)
		return
	}

	//现在key是存在的
	if len(getResp.Kvs) != 0 {
		fmt.Println("当前值：", string(getResp.Kvs[0].Value))
	}

	//从哪里开始监听
	//要从当前的事务ID的下一个事务ID开始监听
	//当前ETCD集群事务I，单调递增的
	watchStartRevision = getResp.Header.Revision + 1

	//我们知道从哪里开始监听之后，就可以创建监听器了
	//创建一个watcher
	watcher = clientv3.NewWatcher(client)

	//启动监听
	fmt.Println("从该版本向后监听:", watchStartRevision)

	//创建一个可以取消的上下文
	ctx, cancelFunc := context.WithCancel(context.TODO())
	//五秒后取消掉
	time.AfterFunc(5 * time.Second, func() {
		//取消ctx
		cancelFunc()
	})

	//传入一个监听的起始事务id
	watchRespChan = watcher.Watch(ctx, "/cron/jobs/job8", clientv3.WithRev(watchStartRevision))

	//处理kv变化事件
	for watchResp = range watchRespChan {
		//events数组可能是不同key的变化，因为它是打包传递过来的
		for _, event = range watchResp.Events{
			switch event.Type {
			case mvccpb.PUT:
				fmt.Println("修改为:", string(event.Kv.Value), "Revision:", event.Kv.CreateRevision, event.Kv.ModRevision)
			case mvccpb.DELETE:
				//这次删除会带来一个修改版本
				fmt.Println("删除了","Revision:", event.Kv.ModRevision)

				
			}
		}
	}

	
}

func etcdTest7() {
	//租约对分布式乐观锁很重要
	var (
		config clientv3.Config
		client *clientv3.Client
		err error
		lease clientv3.Lease
		leaseGrantResp *clientv3.LeaseGrantResponse
		leaseId clientv3.LeaseID
		putResp *clientv3.PutResponse
		getResp *clientv3.GetResponse
		keepResp *clientv3.LeaseKeepAliveResponse
		//这是一个只读channel,写入是有etcd的sdk来做的
		keepRespChan <-chan *clientv3.LeaseKeepAliveResponse
		kv clientv3.KV
	)

	config = clientv3.Config{
		Endpoints:[]string{"127.0.0.1:2379"},
		DialTimeout:5 * time.Second,
	}

	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}

	//申请一个lease（租约）
	lease = clientv3.NewLease(client)

	//grant是申请租约 keepalive是续租,grant是秒为粒度的,申请一个10秒的租约
	//设置超时时间是为了避免程序挂了后还继续占用着锁
	//但程序每down掉，应该让它保持锁，因此需要租约的续租
	if leaseGrantResp, err = lease.Grant(context.TODO(), 10); err != nil {
		fmt.Println(err)
		return
	}

	//如果成功，put一个kv，让它与租约关联起来，从而实现10秒后自动过期

	//拿到租约ID
	leaseId = leaseGrantResp.ID


	//定义一个五秒后会自动过期的context
	//ctx, _ := context.WithTimeout(context.TODO(), 5 * time.Second)

	//续租了5秒，停止续租，10秒的生命期 ＝ 15秒的生命期
	//停止续租，并不代表生命期完了，仍旧有10秒的生命期

	//keepalive会启动一个协程序，自动续租,它的response会定期的回一个
	if keepRespChan, err = lease.KeepAlive(context.TODO(), leaseId); err != nil {
		fmt.Println(err)
		return
	}

	//启动一个协程，去消费keepalive里面的应答
	//处理续租应答的协程
	go func() {
		for {
			select {
			case keepResp = <- keepRespChan:
				//如果续租的过程种出现了异常，也就是sdk和etcd失联了很久
				if keepRespChan == nil {
					fmt.Println("租约已经失效了")
					//结束整个协程
					goto END
				} else{
					//每秒会续租一次，所以就会收到一次应答
					fmt.Println("续租成功了,收到自动续租应答：", keepResp.ID)
				}
			}
		}
		END:
	}()

	//获取kv api子集
	kv = clientv3.NewKV(client)

	//key和lease建立关联
	if putResp, err = kv.Put(context.TODO(), "/cron/lock/job1", "", clientv3.WithLease(leaseId)); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("写入成功：", putResp.Header.Revision)

	//定时的看一下key过期了没有
	for {
		if getResp, err = kv.Get(context.TODO(), "/cron/lock/job1"); err!=nil {
			fmt.Println(err)
			return
		}
		//如果过期了，getResp 肯定是0
		if getResp.Count == 0 {
			fmt.Println("kv过期了")
			break
		}

		//否则还没过期
		fmt.Println("还没过期:", getResp.Kvs)
		time.Sleep(2 * time.Second)
	}


}

func etcdTest6() {
	//删除
	var (
		config clientv3.Config
		client *clientv3.Client
		err error
		kv clientv3.KV
		delResp *clientv3.DeleteResponse
		kvpair *mvccpb.KeyValue
	)

	config = clientv3.Config{
		Endpoints:[]string{"127.0.0.1:2379"},
		DialTimeout:5 * time.Second,
	}

	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}

	kv = clientv3.NewKV(client)

	//想获取删除之前是什么值
	//if delResp, err = kv.Delete(context.TODO(), "/cron/jobs/job2", clientv3.WithPrevKV()); err != nil {
	//	fmt.Println(err)
	//	return
	//}

	//删除一个目录下的所有key
	//if delResp, err = kv.Delete(context.TODO(), "/cron/jobs/", clientv3.WithPrefix()); err != nil {
	//	fmt.Println(err)
	//	return
	//}

	//删除从某一个key开始的若干个key,limit只删除两个
	if delResp, err = kv.Delete(context.TODO(), "/cron/jobs/job1", clientv3.WithFromKey(), clientv3.WithLimit(2)); err!= nil {

	}

	//如果被删除有元素，deleted就不会是零
	//也可以用len(delResp.Prekvs)判断

	//被删除之前的value是什么
	if len(delResp.PrevKvs) != 0 {
		for _, kvpair = range delResp.PrevKvs {
			fmt.Println("删除了：", string(kvpair.Key), string(kvpair.Value))
		}
	}
}


func etcdTest5() {
	//按目录获取
	var (
		config clientv3.Config
		client *clientv3.Client
		err error
		kv clientv3.KV
		getResp *clientv3.GetResponse
	)

	config = clientv3.Config{
		Endpoints:[]string{"127.0.0.1:2379"},
		DialTimeout:5 * time.Second,
	}

	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}

	kv = clientv3.NewKV(client)

	//读取/cron/jobs/为前缀的所有key
	//有了withprefix，就可以找到key为前缀的所有key
	if getResp, err = kv.Get(context.TODO(), "/cron/jobs/", clientv3.WithPrefix()); err != nil {
		fmt.Println(err)
	} else {
		//获取成功，遍历所有的kvs
		fmt.Println(getResp.Kvs)
	}
}

func etcdTest4() {
	//get
	var (
		config clientv3.Config
		client *clientv3.Client
		err error
		kv clientv3.KV
		getResp *clientv3.GetResponse
	)

	config = clientv3.Config{
		Endpoints:[]string{"127.0.0.1:2379"},
		DialTimeout:5 * time.Second,
	}

	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}

	//用于读写etcd的键值对
	kv = clientv3.NewKV(client)

	if getResp, err = kv.Get(context.TODO(), "/cron/jobs/job1", /*clientv3.WithCountOnly()*/); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(getResp.Kvs, getResp.Count)
	}
}

func etcdTest3(){
	//put 操作
	var (
		config clientv3.Config
		client *clientv3.Client
		err error
		kv clientv3.KV
		putResp *clientv3.PutResponse
	)

	config = clientv3.Config{
		Endpoints:[]string{"127.0.0.1:2379"},
		DialTimeout:5 * time.Second,
	}

	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}

	//用于读写etcd的键值对
	kv = clientv3.NewKV(client)

	//context可以用于超时关闭，如果想什么都不做，就context.TO DO()，这样就不会超时关闭
	if putResp, err = kv.Put(context.TODO(), "/cron/jobs/job1", "hello", clientv3.WithPrevKV()); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Revision:", putResp.Header.Revision)
		//如果非空，说明覆盖了原来某个值
		if putResp.PrevKv != nil {
			fmt.Println("PrevValue:", string(putResp.PrevKv.Value))
		}
	}

}

func etcdTest2()  {
	var (
		config clientv3.Config
		client *clientv3.Client
		err error
	)

	config = clientv3.Config{
		Endpoints:[]string{"127.0.0.1:2379"},
		DialTimeout:5 * time.Second,
	}

	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}

	client = client
}

func etcdTest1() {

	var (
		dialTimeout = 5 * time.Second
		requestTimeout = 2 * time.Second
		endpoints = []string{"127.0.0.1:2379"}
	)

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:endpoints,
		DialTimeout:dialTimeout,
	})

	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	log.Println("存储值")
	if _, err := cli.Put(context.TODO(), "/sensors", `{sensor01:{topic:"w_sensor01"}}`); err != nil {
		log.Fatal(err)
	}

	log.Println("获取值")
	if resp, err := cli.Get(context.TODO(), "/sensors"); err != nil {
		log.Fatal(err)
	} else {
		log.Println("resp: ", resp)
	}

	log.Println("事务&超时")
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	_, err = cli.Txn(ctx).
		If(clientv3.Compare(clientv3.Value("key"), ">", "abc")).
		Then(clientv3.OpPut("key", "XYZ")).
		Else(clientv3.OpPut("key", "ABC")).
		Commit()
	cancel()
	if err != nil {
		fmt.Println(err)
	}

	log.Println("获取tx值")
	if resp, err := cli.Get(context.TODO(), "key"); err != nil {
		log.Fatal(err)
	} else {
		log.Println("resp: ", resp)
	}


	log.Println("watch..")
	rch := cli.Watch(context.Background(), "/test/hello", clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
	}

	if err != nil {
		println(err)
	}

}
