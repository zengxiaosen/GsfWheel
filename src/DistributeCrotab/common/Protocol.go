package common

import (
	"context"
	"encoding/json"
	"github.com/gorhill/cronexpr"
	"strings"
	"time"
)

//定时任务
type Job struct {
	Name     string `json:"name"`     //任务名
	Command  string `json:"command"`  //shell命令
	CronExpr string `json:"cronExpr"` //cron表达式
}

//任务调度计划
type JobSchedulePlan struct {
	Job      *Job                 //要调度的任务信息
	Expr     *cronexpr.Expression //解析好的cronexpr表达式
	NextTime time.Time            //下次调度时间
}

//任务执行状态
type JobExecuteInfo struct {
	//任务信息
	Job *Job
	//理论上的调度时间
	PlanTime time.Time
	//实际的调度时间
	RealTime time.Time
	//任务command的context
	CancelCtx context.Context
	//用于取消command执行的cancel函数
	CancelFunc context.CancelFunc
}

//HTTP接口应答
type Response struct {
	Errno int         `json:"errno"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}

//变化事件
type JobEvent struct {
	EventType int //SAVE, DELETE
	//变化的信息
	Job *Job
}

//任务执行结果
type JobExecuteResult struct {
	//执行状态
	ExecuteInfo *JobExecuteInfo
	//脚本输出
	Output []byte
	//脚本错误原因
	Err error
	//启动时间(脚本真实的启动时间）
	StartTime time.Time
	//结束时间
	EndTime time.Time
}

//任务执行日志
type JobLog struct {
	JobName      string `json:"jobName" bson:"jobName"`           //任务名字
	Command      string `json:"command" bson:"command"`           //脚本命令
	Err          string `json:"err" bson:"err"`                   //错误原因
	Output       string `json:"output" bson:"output"`             //脚本输出
	PlanTime     int64  `json:"planTime" bson:"planTime"`         //计划开始时间
	ScheduleTime int64  `json:"scheduleTime" bson:"scheduleTime"` //实际调度时间
	StartTime    int64  `json:"startTime" bson:"startTime"`       //任务执行开始时间
	EndTime      int64  `json:"endTime" bson:"endTime"`           //任务执行结束时间
}

//任务日志过滤条件
type JobLogFilter struct {
	JobName string `bson:"jobName"`
}

//任务日志排序条件
type SortLogByStartTime struct {
	SortOrder int `bson:"startTime"` //{startTime:-1}
}

//日志批次
type LogBatch struct {
	Logs []interface{} //多条日志
}

func BuildResponse(errno int, msg string, data interface{}) (resp []byte, err error) {
	var (
		response Response
	)

	response.Errno = errno
	response.Msg = msg
	response.Data = data

	resp, err = json.Marshal(response)
	return
}

//反序列化job
func UnpackJob(value []byte) (ret *Job, err error) {

	var (
		job *Job
	)

	job = &Job{}
	if err = json.Unmarshal(value, job); err != nil {
		return
	}

	ret = job

	return
}

//从etcd的key中提取任务名
func ExtractKillerName(jobKey string) string {
	return strings.TrimPrefix(jobKey, JOB_SAVE_DIR)
}

//从 /cron/killer/job10提取job10
func ExtractJobName(killerKey string) string {
	return strings.TrimPrefix(killerKey, JOB_KILLER_DIR)
}

//任务变化事件有2种： 1）更新任务  2）删除任务
func BuildJobEvent(eventType int, job *Job) (jobEvent *JobEvent) {
	return &JobEvent{
		EventType: eventType,
		Job:       job,
	}
}

//构造任务执行计划
func BuildJobSchedulePlan(job *Job) (jobSchedulelan *JobSchedulePlan, err error) {
	var (
		expr *cronexpr.Expression
	)

	//解析JOB的cron表达式
	if expr, err = cronexpr.Parse(job.CronExpr); err != nil {
		return
	}

	//生成任务调度计划对象
	jobSchedulelan = &JobSchedulePlan{
		Job:  job,
		Expr: expr,
		//expr表达式，传入当前时间，就可以得到下次执行时间
		NextTime: expr.Next(time.Now()),
	}

	return
}

//构造执行状态信息,计划过期了，要生成一个执行状态，执行起来
func BuildJobExecuteInfo(jobSchedulePlan *JobSchedulePlan) (jobExecuteInfo *JobExecuteInfo) {
	jobExecuteInfo = &JobExecuteInfo{
		Job: jobSchedulePlan.Job,
		//计划调度时间
		PlanTime: jobSchedulePlan.NextTime,
		//真是调度时间
		RealTime: time.Now(),
	}
	//创建一个可以用于取消的上下文
	jobExecuteInfo.CancelCtx, jobExecuteInfo.CancelFunc = context.WithCancel(context.TODO())
	return
}

//提取worker的ip
func ExtractWorkerIp(regKey string) string {
	return strings.TrimPrefix(regKey, JOB_WORKER_DIR)
}
