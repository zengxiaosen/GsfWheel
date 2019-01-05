package main

import (
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	pb "gsf/src/QosBreaker/helloworld"
	"log"
	"os"
	"time"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
	serviceName = "qosCircuitBreakerService"
)

func main() {

	conn, err := grpc.Dial(address,

		grpc.WithInsecure(),

		//开启grpc中间件的重试功能
		grpc.WithUnaryInterceptor(

			grpc_retry.UnaryClientInterceptor(
				//重试间隔时间
				grpc_retry.WithBackoff(grpc_retry.BackoffLinear(time.Duration(1)*time.Millisecond)),
				//重试次数
				grpc_retry.WithMax(3),
				//重试时间
				grpc_retry.WithPerRetryTimeout(time.Duration(5)*time.Millisecond),
				//返回码如下时重试
				grpc_retry.WithCodes(codes.ResourceExhausted, codes.Unavailable, codes.DeadlineExceeded),
			),


		),
	)

	//熔断器
	hystrix.ConfigureCommand(
		//熔断器名字，可以用服务名称命名，一个名字对应一个熔断器，对应一份熔断策略
		serviceName,
		hystrix.CommandConfig{
			//超时时间 100ms
			Timeout: 100,
			//最大并发数，超过并发返回错误
			MaxConcurrentRequests: 2,
			//请求数量的阀值，用这些数量的请求来计算阈值
			RequestVolumeThreshold: 4,
			//错误率阈值，达到阈值，启动熔断器，25%
			ErrorPercentThreshold: 25,
			//熔断尝试恢复时间，1000ms
			SleepWindow: 1000,
		},
	)

	// conn, err := grpc.Dial(address, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	//test breaker ------------------------------------------------------------------------------------------
	//for i := int64(1); i < 1000; i++ {
	//	for j := int64(1); j < 10; j++ {
	//
	//		//{
	//		//	/*
	//		//	熔断，阻塞方式调用
	//		//	*/
	//		//	req := &pb.HelloRequest{Name: name}
	//		//	var res *pb.HelloReply
	//		//	breakererr := hystrix.Do(serviceName,
	//		//
	//		//		func() error {
	//		//			//正常业务逻辑，一般时访问其他静态资源
	//		//			var berr error
	//		//			//设置总体超时时间 10ms 超时
	//		//			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10*time.Millisecond))
	//		//			defer cancel()
	//		//			res, berr = client.SayHello(
	//		//				ctx, req,
	//		//				// 这里可以再次设置重试次数，重试时间，重试返回码
	//		//				grpc_retry.WithMax(3),
	//		//				grpc_retry.WithPerRetryTimeout(time.Duration(5)*time.Millisecond),
	//		//				grpc_retry.WithCodes(codes.DeadlineExceeded),
	//		//			)
	//		//			return berr
	//		//		},
	//		//
	//		//		func(err error) error {
	//		//
	//		//			/*
	//		//			失败处理逻辑，访问其他资源失败时，或者处于熔断开启状态时，会调用这段逻辑
	//		//			可以简单构造一个response返回，也可以有一定的策略，比如访问备份资源
	//		//			也可以直接返回err，这样不用和远端失败的资源通信，防止雪崩
	//		//			这里简单返回一个response
	//		//			 */
	//		//			fmt.Println(err)
	//		//			res = &pb.HelloReply{Message:"阻塞方式 failback: service breaker err response"}
	//		//			return nil
	//		//		},
	//		//
	//		//	)
	//		//	if breakererr != nil {
	//		//		//事实上这个断言永远为假，因为错误会触发熔断调用 fallback，而 fallback 函数返回 nil
	//		//		fmt.Printf("sent req to server failed. err: [%v]\n", err)
	//		//	}
	//		//	fmt.Println(req, res)
	//		//	log.Printf("get server response: %s", res.Message)
	//		//
	//		//
	//		//}
	//
	//
	//
	//
	//		{
	//			/*
	//			熔断，非阻塞方式调用
	//			建议在有多个资源需要并发访问的场景下是使用
	//			*/
	//			req := &pb.HelloRequest{Name: name}
	//			var res1 *pb.HelloReply
	//			var res2 *pb.HelloReply
	//
	//			success := make(chan struct{}, 2)
	//
	//			//无 fallback 处理
	//			errc1 := hystrix.Go(serviceName, func() error {
	//				var err error
	//				ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10*time.Millisecond))
	//				defer cancel()
	//				res1, err = client.SayHello(ctx, req)
	//				if err == nil {
	//					success <- struct{}{}
	//				}
	//				return err
	//			}, nil)
	//
	//			//有 fallback 处理
	//			errc2 := hystrix.Go(serviceName, func() error {
	//				var err error
	//				ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10*time.Millisecond))
	//				defer cancel()
	//				res2, err = client.SayHello(ctx, req)
	//				if err == nil {
	//					success <- struct{}{}
	//				}
	//				return err
	//			}, func(err error) error {
	//				fmt.Println(err)
	//				res2 = &pb.HelloReply{Message:"非阻塞方式 failback: service breaker err response"}
	//				fmt.Println("非阻塞方式 failback: service breaker err response")
	//				success <- struct{}{}
	//				return nil
	//			})
	//
	//			for i := 0; i< 2; i++ {
	//				select {
	//				case <-success:
	//					fmt.Println("success", i)
	//				case err := <- errc1:
	//					fmt.Println("err1: ", err)
	//				case err := <- errc2:
	//					//这个分之永远不会走，因为熔断机制里面永远不会返回错误
	//					fmt.Println("err2:", err)
	//
	//				}
	//			}
	//
	//			fmt.Println("req: ", req, ", res1: ", res1, ", res2: ", res2)
	//
	//		}
	//	}
	//
	//}


	//------------------------------------------------------------------------------------------




	{
		/*
		熔断，阻塞方式调用
		*/
		req := &pb.HelloRequest{Name: name}
		var res *pb.HelloReply
		breakererr := hystrix.Do(serviceName,

			func() error {
				//正常业务逻辑，一般时访问其他静态资源
				var berr error
				//设置总体超时时间 10ms 超时
				ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10*time.Millisecond))
				defer cancel()
				res, berr = client.SayHello(
					ctx, req,
					// 这里可以再次设置重试次数，重试时间，重试返回码
					grpc_retry.WithMax(3),
					grpc_retry.WithPerRetryTimeout(time.Duration(5)*time.Millisecond),
					grpc_retry.WithCodes(codes.DeadlineExceeded),
				)
				return berr
			},

			func(err error) error {

				/*
				失败处理逻辑，访问其他资源失败时，或者处于熔断开启状态时，会调用这段逻辑
				可以简单构造一个response返回，也可以有一定的策略，比如访问备份资源
				也可以直接返回err，这样不用和远端失败的资源通信，防止雪崩
				这里简单返回一个response
				 */
				fmt.Println(err)
				res = &pb.HelloReply{Message:"阻塞方式 failback: service breaker err response"}
				return nil
			},

		)
		if breakererr != nil {
			//事实上这个断言永远为假，因为错误会触发熔断调用 fallback，而 fallback 函数返回 nil
			fmt.Printf("sent req to server failed. err: [%v]\n", err)
		}
		fmt.Println(req, res)
		log.Printf("get server response: %s", res.Message)


	}

	{
		/*
		熔断，非阻塞方式调用
		建议在有多个资源需要并发访问的场景下是使用
		*/
		req := &pb.HelloRequest{Name: name}
		var res1 *pb.HelloReply
		var res2 *pb.HelloReply

		success := make(chan struct{}, 2)

		//无 fallback 处理
		errc1 := hystrix.Go(serviceName, func() error {
			var err error
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10*time.Millisecond))
			defer cancel()
			res1, err = client.SayHello(ctx, req)
			if err == nil {
				success <- struct{}{}
			}
			return err
		}, nil)

		//有 fallback 处理
		errc2 := hystrix.Go(serviceName, func() error {
			var err error
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10*time.Millisecond))
			defer cancel()
			res2, err = client.SayHello(ctx, req)
			if err == nil {
				success <- struct{}{}
			}
			return err
		}, func(err error) error {
			fmt.Println(err)
			res2 = &pb.HelloReply{Message:"非阻塞方式 failback: service breaker err response"}
			success <- struct{}{}
			return nil
		})

		for i := 0; i< 2; i++ {
			select {
			case <-success:
				fmt.Println("success", i)
			case err := <- errc1:
				fmt.Println("err1: ", err)
			case err := <- errc2:
				//这个分之永远不会走，因为熔断机制里面永远不会返回错误
				fmt.Println("err2:", err)

			}
		}

		fmt.Println("req: ", req, ", res1: ", res1, ", res2: ", res2)

	}


}
