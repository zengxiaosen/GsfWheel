package qosRateLimit

import (
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"strconv"
	"sync"
)
var(
	lock sync.RWMutex
)
func LimitRateServerInterceptor() grpc.UnaryServerInterceptor {


	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		//step1

		stat := true
		//一个rpc请求，因为目前不做批处理
		querySize := 1

		//初始化令牌桶容量 暂时不给用户自定义链
		//capacity = bucketCapacity
		log.Println("int64(querySizePerSecond): ", int64(querySize))
		lock.Lock()
		ret := grantToken(int64(querySize));
		lock.Unlock()
		if ret {
			fmt.Println("trans packet")
			stat = true
		} else {
			fmt.Println("No trans")
			stat = false

		}

		//step2

		if stat == true {
			log.Println("stat is true")
			return handler(ctx, req)
		} else {
			log.Println("stat is false")
			return nil, errors.New("服务并发限制, "+ "当前服务最大并发数为" + strconv.FormatInt(capacity, 10) + ", 请稍后重试")
		}
	}
}
