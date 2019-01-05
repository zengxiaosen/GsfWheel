package main

import (
	"flag"
	"fmt"
	"gsf/src/DistributeCrotab/worker"
	"runtime"
	"time"
)

var (
	//配置文件路径,通过flag去解析
	confFile string
)

// 解析命令行参数
func initArgs() {
	// worker -config ./worker.json
	// worker -h
	flag.StringVar(&confFile, "config", "/Users/zeng/go/src/gsf/src/DistributeCrotab/worker/main/worker.json", "worker.json")
	flag.Parse()
}

// 初始化线程数量
func initEnv() {
	//设置线程数和cpu的数量相等
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var (
		err error
	)

	//初始化命令行参数
	initArgs()

	//初始化线程
	initEnv()

	//加载配置
	if err = worker.InitConfig(confFile); err != nil {
		goto ERR
	}

	//服务注册
	if err = worker.InitRegister(); err != nil {
		goto ERR
	}

	//启动日志协程
	if err = worker.InitLogSink(); err != nil {
		goto ERR
	}

	//启动执行器
	if err = worker.InitExecutor(); err != nil {
		goto ERR
	}
	//启动调度器
	if err = worker.InitScheduler(); err != nil {
		goto ERR
	}

	//初始化任务管理器
	if err = worker.InitJobMgr(); err != nil {
		goto ERR
	}

	//正常退出
	for {
		time.Sleep(1 * time.Second)
	}
	return

ERR:
	fmt.Println(err)
}
