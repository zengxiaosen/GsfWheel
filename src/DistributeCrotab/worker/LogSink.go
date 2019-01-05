package worker

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gsf/src/DistributeCrotab/common"
	"time"
)

//mongodb存储日志
type LogSink struct {
	logCollection  *mgo.Collection
	logChan        chan *common.JobLog //scheduler把日志扔到这里存储
	autoCommitChan chan *common.LogBatch
}

var (
	G_logSink *LogSink
)

func (logSink *LogSink) saveLogs(batch *common.LogBatch) {

	var (
		b []bson.M
		//logs []interface{}
	)
	//TODO mongo的批量插入
	for _, m := range batch.Logs {

		data := m.(*common.JobLog)
		body := *data
		fmt.Println(body)
		bodyByte, err := bson.Marshal(body)
		if err != nil {
			panic(err)
		}
		mmap := bson.M{}
		err = bson.Unmarshal(bodyByte, mmap)
		if err != nil {
			panic(err)
		}

		b = append(b, mmap)
		err = logSink.logCollection.Insert(mmap)
		if err != nil {
			fmt.Println(err)
		}
	}

}

//日志存储协程
func (logSink *LogSink) writeLoop() {
	var (
		log      *common.JobLog
		logBatch *common.LogBatch //当前的批次
		//让这个批次超时自动提交，给它1秒的时间
		commitTimer  *time.Timer
		timeoutBatch *common.LogBatch //超时批次
	)

	for {
		select {
		case log = <-logSink.logChan:
			//把这条log写到mongodb中
			//logSink.logCollection.insertone
			//每次插入需要等待mongodb的一次请求往返，耗时可能因为网络慢话费较长的时间

			if logBatch == nil {
				logBatch = &common.LogBatch{}

				//让这个批次超时自动提交（给1秒的时间）
				commitTimer = time.AfterFunc(
					time.Duration(G_config.JobLogCommitTimeout)*time.Millisecond, func(batch *common.LogBatch) func() {
						return func() {
							logSink.autoCommitChan <- batch
						}
					}(logBatch),
				)
			}

			//把新日志追加到批次中
			logBatch.Logs = append(logBatch.Logs, log)

			//如果批次满了，就立即发送
			if len(logBatch.Logs) >= G_config.JobLogBatchSize {
				//发送日志
				logSink.saveLogs(logBatch)
				//晴空logBatch
				logBatch = nil
				//取消定时器
				commitTimer.Stop()
			}
		case timeoutBatch = <-logSink.autoCommitChan: //过期的批次

			//判断过期批次是否仍旧是当前批次
			if timeoutBatch != logBatch {
				//证明logBatch已经被清空了，升至进入下一批次的batch，不等于证明这个过期的batch已经被提交过了
				continue // 跳过已被提交的批次
			}

			//把批次写入到mongo中
			logSink.saveLogs(timeoutBatch)
			//清空logBatch
			logBatch = nil

		}
	}
}

func InitLogSink() (err error) {

	session, err := mgo.Dial(G_config.MongodbUri)
	if err != nil {
		panic(err)
	}
	//defer session.Close()

	G_logSink = &LogSink{
		logCollection:  session.DB("test").C("log"),
		logChan:        make(chan *common.JobLog, 1000),
		autoCommitChan: make(chan *common.LogBatch, 1000),
	}

	//启动一个mongodb处理协程消费队列
	go G_logSink.writeLoop()

	return
}

//发送日志
func (logSink *LogSink) Append(jobLog *common.JobLog) {
	select {
	case logSink.logChan <- jobLog:

	default:
		//队列满了就丢弃

	}
}
