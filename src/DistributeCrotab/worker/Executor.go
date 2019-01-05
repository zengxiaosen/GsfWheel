package worker

import (
	"gsf/src/DistributeCrotab/common"
	"math/rand"
	"os/exec"
	"time"
)

//任务执行器
type Executor struct {
}

var (
	G_executor *Executor
)

//执行一个任务
func (executor *Executor) ExecuteJob(info *common.JobExecuteInfo) {
	go func() {
		var (
			cmd     *exec.Cmd
			err     error
			output  []byte
			result  *common.JobExecuteResult
			jobLock *JobLock
		)

		//定义任务结果
		result = &common.JobExecuteResult{
			ExecuteInfo: info,
			Output:      make([]byte, 0),
		}

		//初始化锁
		jobLock = G_jobMgr.CreateJobLock(info.Job.Name)

		//纪录任务开始时间
		result.StartTime = time.Now()

		//上锁
		//随机睡眠(0~1s)
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		err = jobLock.TryLock()
		//释放锁
		defer jobLock.Unlock()

		if err != nil {
			//上锁失败
			result.Err = err
			result.EndTime = time.Now()

		} else {

			//上锁成功后，重置任务启动时间
			result.StartTime = time.Now()

			//执行shell命令
			cmd = exec.CommandContext(info.CancelCtx, "/bin/bash", "-c", info.Job.Command)

			//执行并捕获输出
			//如果任务被调度器强杀了，这里的err会返回一个错误
			output, err = cmd.CombinedOutput()

			//纪录任务结束时间
			result.EndTime = time.Now()
			result.Output = output
			result.Err = err

			//任务执行完成后，把执行的结果返回给Scheduler，Scheduler会从executingTable中删除掉执行纪录

		}

		//完成从执行器向scheduler的执行通知
		G_scheduler.PushJobResult(result)

	}()
}

//初始化执行器
func InitExecutor() (err error) {
	G_executor = &Executor{}
	return
}
