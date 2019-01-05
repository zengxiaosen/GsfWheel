package master

import (
	"gopkg.in/mgo.v2"
	"gsf/src/DistributeCrotab/common"
)

var (
	//单例
	G_logMgr *LogMgr
)

//mongodb存储日志
type LogMgr struct {
	logCollection *mgo.Collection
}

func InitLogMgr() (err error) {
	session, err := mgo.Dial(G_config.MongodbUri)
	if err != nil {
		panic(err)
	}
	//defer session.Close()

	G_logMgr = &LogMgr{
		logCollection: session.DB("test").C("log"),
	}

	return

}

//查看任务日志
func (logMgr *LogMgr) ListLog(name string, skip int, limit int) ([]common.JobLog, error) {

	var (
		filter     *common.JobLogFilter
		jobLogInfo *[]common.JobLog
		logArr     []common.JobLog
	)

	logArr = make([]common.JobLog, 0)

	//过滤条件
	filter = &common.JobLogFilter{JobName: name}

	//info := []common.JobLog
	jobLogInfo = &[]common.JobLog{}
	logMgr.logCollection.Find(filter).Sort("-startTime").Skip(skip).Limit(limit).All(jobLogInfo)
	for _, info := range *jobLogInfo {
		logArr = append(logArr, info)
	}

	return logArr, nil

}
