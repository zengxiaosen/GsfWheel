package main

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

func main() {

	//opTest1()

	opTest2()

}

type result struct {
	err error
	output []byte
}

func opTest2() {
	//执行1个cmd，让它在一个协程里去执行，让它执行2秒：sleep 2； echo hello
	//1秒的时候，我们杀死cmd

	var (
		ctx context.Context
		cancelFunc context.CancelFunc
		cmd *exec.Cmd
		resultChan chan *result
		res *result
	)

	//创建一个结果队列
	resultChan = make(chan *result, 1000)

	//context: chan byte
	//cancelFunc: close(chan byte)

	//withCancel的这个上下文继承了todo这个context，是继承的关系
	ctx, cancelFunc = context.WithCancel(context.TODO())

	go func() {

		var(
			output []byte
			err error
		)
		//在写成里面要去调command，是没办法被取消的，写成里面应该调commandWithContext

		/*
		param1 context golang里面的标准库
		这个context方法是用来取消命令执行的
		CommandContext方法内部回有个select {case <- ctx.Done();}去监听这个ctx是否被关闭，它会实时感知到
		一旦我们主进程执行了cancelFunc()，就会被感知到
		感知到之后，CommandContext就会哟过kill pid,进程ID，把我们系统bash程序杀死
		 */
		cmd = exec.CommandContext(ctx, "/bin/bash", "-c", "sleep 2;echo hello;")
		//执行任务，捕获输出
		output, err = cmd.CombinedOutput()

		//这里要传到主协程里面，两个协程之间传递用channel
		resultChan <- &result{
			err: err,
			output:output,
		}

	}()

	//继续往下走
	time.Sleep(1 * time.Second)

	//取消上下文
	cancelFunc()

	//在main协程里，等待子协程的退出，并打印任务执行结果
	res = <- resultChan

	//打印任务执行结果
	fmt.Println(res.err, string(res.output))
}

func opTest1() {
	var (
		cmd *exec.Cmd
		err error
		output []byte
	)

	//生成cmd
	cmd = exec.Command("/bin/bash", "-c", "sleep 5;ls -l;sleep 1;echo hello;")

	//创建子进程
	//err = cmd.Run()
	//fmt.Println(err)

	//执行子进程，并捕获输出，捕获方法是利用linux的pipe
	if output, err = cmd.CombinedOutput(); err != nil {
		fmt.Println(err)
		return
	}
	//打印子进程的输出
	fmt.Println(string(output))
}