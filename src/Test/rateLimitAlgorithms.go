package main

import (
	"fmt"
	"time"
)
/*
单机level
 */

func main() {
	//计数器算法法
	//algo1()
	//漏都算法
	//algo2()
	//令牌桶算法
	algo3()
}

/*
计数器算法：

计数器法是限流算法里最简单也是最容易实现的一种算法。比如我们规定，对于A接口来说，我们1分钟的访问次数不能超过100个。那么我们可以这么做：
在一开始的时候，我们可以设置一个计数器counter，每当一个请求过来的时候，counter就加1，如果counter的值大于100并且该请求于第一个请求的
时间间隔还在1分钟之内，那么说明请求数过多，如果该请求于第一个请求的时间间隔大于1分钟，且counter的值还在限流范围内，那么就重置counter

缺点：
这种实现方式，相信大家都知道有一个弊端：如果我在单位时间1s内的前10ms，已经通过了100个请求，那后面的990ms，只能眼巴巴的把请求拒绝，
我们把这种现象称为“突刺现象”

 */

func algo1()  {

	timeStamp := time.Now()
	fmt.Println("timeStamp: ", timeStamp)
	as := algo1Context{0, 100, 1000}

	grantFlag := grant(&as)
	fmt.Println("is it possible to grant: ", grantFlag)



}

/*
漏斗算法：

漏斗有一个进水口和一个出水口，出水口以一定速率出水，并且有一个最大出水速率：

在漏斗中没有水的时候：
1，如果进水速率小于等于最大出水速率，那么，出水速率等于进水速率，此时，不会积水
2，如果进水速率大于最大出水速率，那么，漏斗以最大速率出水，此时，多余的水会积在漏斗中

在漏斗中有水的时候：
1，出水口以最大速率出水
2，如果漏斗未满，且有进水的话，那么这些水会积在漏斗中
3，如果漏斗已满，且有进水的话，那么这些水会溢出到漏斗之外


*/
func algo2() {

}

/*
令牌桶算法：网络流量整形和速率限制中最常用的一种算法。

对于很多应用场景来说，除了要求能够限制数据的平均传输速率外，还要求允许某种程度的突发传输。这时候漏桶算法可能就不适合了，令牌桶算法更适合。

令牌桶算法的原理是系统以恒定的速率产生令牌，然后把令牌放到令牌桶中，令牌桶有一个容量，当令牌桶满了的时候，再向其中放令牌，
那么多余的令牌会被丢弃；当想要处理一个请求的时候，需要从令牌桶中取出一个令牌，如果此时令牌桶中没有令牌，那么则拒绝该请求。

 */

func algo3() {

	if ret := grant3(1024 * 1024 * 128); ret {
		fmt.Println("trans packet at first")
		return
	}

	select {
	case <-time.After(time.Second):
		break
	}

	if ret := grant3(1024 * 1024 * 128); ret {
		fmt.Println("trans packet at second")
		return
	}
	fmt.Println("No trans")


}

//packetSize 此次请求分配的令牌数量
func grant3 (packetsize int64) bool {
	if packetsize > capacity || packetsize < 1 {
		return false
	}

	now := time.Now()
	//获取令牌桶的令牌数量---如果生成的令牌数量多于桶容量，则溢出
	tokens = min(capacity, tokens+(now.Unix()-timestampglobal.Unix())*rate)
	timestampglobal = now
	if tokens < packetsize {
		return false
	} else {
		tokens -= packetsize
		return true
	}

}

var timestampglobal time.Time = time.Now()

func min(va, vb int64) int64 {
	if va <= vb {
		return va
	}
	return vb
}



//桶容量
var capacity int64 = 256 * 1024 * 1024

//令牌放入的速度
var rate int64 = 128 * 1024 * 1024

//当前令牌数量
var tokens int64 = 0




func grant(as *algo1Context) bool {
	//模拟停留1s
	time.Sleep(time.Second)
	nowTime := time.Now()
	fmt.Println("nowTime: ", nowTime)
	duration := nowTime.Second() - timeStamp.Second()
	fmt.Println("duration(s): ", duration)

	if nowTime.Second() < timeStamp.Second() + as.interval / 1000 {
		//在时间窗口之内
		as.reqCount ++
		//判断当前时间窗口是否超过最大请求控制数
		return as.reqCount <= as.limit
	} else {
		timeStamp = time.Now()
		//超时后重置
		as.reqCount = 1
		return true

	}
}

var (
	timeStamp = time.Now()
)
type algo1Context struct {

	reqCount int
	//时间窗口内最大请求数
	limit int
	//时间窗口ms
	interval int
}