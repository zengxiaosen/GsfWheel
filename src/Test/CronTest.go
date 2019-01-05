package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)

//
//cron基本格式
//1) 分钟(0~59)
//2) 小时(0~23)
//3) 日(1~31)
//4) 月(1~12)
//5) 星期(0~6)
//6) Command shell命令
//
//如果1-5中换成 * ，就是"每"的语义，比如每分钟执行一次
//
//case:
//每五分钟执行一次: */5 * * * * echo hello > /tmp/x.log 除了分钟, 都写成星号, 就是每个时间点都去执行一次
//每个小时的第1分钟, 第二分钟 到 第五分钟, 分别执行一次, 共执行五次 1-5 * * * * /usr/bin/python /data/x.py
//每天的10点, 22点整执行1次: 0 20,22 * * * echo bye | tail -1
//
//golang里面有一个开源的Cronexpr库
//传入一个当前时间, 和一个cron表达式, 能够得到下次执行时间是什么, 然后下次去执行就行了.
//



func main() {

	//cron job 调度一个任务
	//cronTest1()

	//cron job 调度多个任务
	cronTest2()
}

//代表一个任务
type CronJob struct {
	expr *cronexpr.Expression
	//每个任务是否过期，实际上是通过它的nextTime来判断的,本质上是通过expr.Next(now)来执行
	nextTime time.Time
}

func cronTest2(){
	//需要有一个调度协程，它定时检查所有的cron任务，谁过期了就执行谁

	var (
		cronJob *CronJob
		expr *cronexpr.Expression
		now time.Time
		//任务表，key是任务的名字，value就是cron job
		scheduleTable map[string]*CronJob
	)

	//初始make不需要传空间，初始化一个大小为0的
	scheduleTable = make(map[string]*CronJob)

	//当前时间
	now = time.Now()
	//1,我们定义两个cronjob
	expr = cronexpr.MustParse("*/5 * * * * * *")
	cronJob = &CronJob{
		expr:expr,
		nextTime:expr.Next(now),
	}
	//把任务注册到调度表
	scheduleTable["job1"] = cronJob


	//第二个任务
	expr = cronexpr.MustParse("*/5 * * * * * *")
	cronJob = &CronJob{
		expr:expr,
		nextTime:expr.Next(now),
	}
	//任务注册到调度表
	scheduleTable["job2"] = cronJob

	//启动一个调度协程
	go func() {
		var (
			jobName string
			cronJob *CronJob
			now time.Time
		)
		//定时检查一下任务调度表
		for {
			now = time.Now()
			for jobName, cronJob = range scheduleTable {
				//判断是否过期
				if cronJob.nextTime.Before(now) || cronJob.nextTime.Equal(now) {
					//启动一个协程，执行这个任务,这里模拟一下，执行这个任务
					go func(jobName string) {
						fmt.Println("执行:", jobName)
					}(jobName)

					//计算下一次调度时间
					cronJob.nextTime = cronJob.expr.Next(now)
					fmt.Println(jobName,"下次执行时间:", cronJob.nextTime)
				}
			}

			//睡眠100毫秒,定义间隔，避免把cpu打满
			select {
			//将在100毫秒可读，返回
			case <- time.NewTimer(100 * time.Millisecond).C:
			}


		}
	}()

	//避免主协程退出，让它睡100秒
	time.Sleep(100 * time.Second)

}



func cronTest1() {

	var (
		expr *cronexpr.Expression
		err error
		now time.Time
		nexTime time.Time
	)

	//每分钟执行一次
	//if expr, err = cronexpr.Parse("* * * * *"); err != nil {
	//	fmt.Println(err)
	//	return
	//}


	//cronexpr支持7位，秒粒度，也支持年的配置（2018-2099），第一位是秒，第七位是年
	//每五秒执行一次
	if expr, err = cronexpr.Parse("*/5 * * * * * *"); err != nil {
		fmt.Println(err)
		return
	}

	//当前时间
	now = time.Now()
	//下次调度时间
	nexTime = expr.Next(now)

	fmt.Println(now, nexTime)

	//等待这个定时器超时,第一个参数时间后，第二个参数回调函数就会执行,  这里会异步一个协程
	time.AfterFunc(nexTime.Sub(now), func() {
		fmt.Println("被调度了:", nexTime)
	})

	//主协程停5秒
	time.Sleep(15 * time.Second)

	expr = expr




}