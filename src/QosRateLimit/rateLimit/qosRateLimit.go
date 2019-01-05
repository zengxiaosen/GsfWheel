package qosRateLimit

import (
	"log"
	"time"
)

//packetSize 此次请求分配的令牌数量
func grantToken(packetsize int64) bool {

	log.Println("capacity: ", capacity)
	log.Println("packetsize: ", packetsize)


	if packetsize > capacity || packetsize < 1 {

		return false
	}

	now := time.Now()
	log.Println("tokens(before 这个阶段是紧接着上一批次完成后最后的令牌数量): ", tokens)
	//获取令牌桶的令牌数量---如果生成的令牌数量多于桶容量，则溢出
	tokens = min(capacity, tokens+(now.Unix()-timestampglobal.Unix())*rate)
	log.Println("now.Unix(): ", now.Unix())
	log.Println("timestampglobal.Unix(): ", timestampglobal.Unix())
	log.Println("tokens(middle 这个阶段增加本批次生成的令牌数量): ", tokens)
	timestampglobal = now
	if tokens < packetsize {
		return false
	} else {
		log.Println("tokens(after 这个阶段消费令牌，拿到令牌取做业务操作): ", tokens)
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
var rate int64 = 10

//当前令牌数量
var tokens int64 = 0

