package main

import (
	"flag"
	"fmt"
	"gsf/src/DistributeCrotab/master"
	"runtime"
	"time"
)

var (
	//配置文件路径,通过flag去解析
	confFile string
)

// 解析命令行参数
func initArgs() {
	// master -config ./master.json -xxx 123 -yyy ddd
	flag.StringVar(&confFile, "config", "/Users/zeng/go/src/gsf/src/DistributeCrotab/master/main/master.json", "指定master.json")
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
	if err = master.InitConfig(confFile); err != nil {
		goto ERR
	}

	//初始化集群管理器,服务发现模块
	if err = master.InitWorkerMgr(); err != nil {
		goto ERR
	}

	//日志管理器
	if err = master.InitLogMgr(); err != nil {
		goto ERR
	}

	//任务管理器
	if err = master.InitJobMgr(); err != nil {
		goto ERR
	}

	//启动Api Http服务
	if err = master.InitApiServer(); err != nil {
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
