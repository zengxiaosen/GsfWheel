package main

import (
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	pb "gsf/src/QosTracer/helloworld"
	"log"
	"os"
	"time"

	"github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
	serviceName = "qosCircuitBreakerService"
)

func main() {

	//zipkin
	collector, err := zipkin.NewHTTPCollector("http://localhost:9411/api/v1/spans")
	if err != nil {
		panic(err)
		return
	}
	defer collector.Close()

	tracer, err := zipkin.NewTracer(
		zipkin.NewRecorder(collector, false, "localhost:0", "grpc_client"),
		zipkin.ClientServerSameSpan(true),
		zipkin.TraceID128Bit(true),
	)
	if err != nil {
		panic(err)
	}
	opentracing.InitGlobalTracer(tracer)

	// Set up a connection to the server.

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

		grpc.WithUnaryInterceptor(
			otgrpc.OpenTracingClientInterceptor(tracer),
		),

	)




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

	// Create Root Span for duration of the interaction with svc1
	span := opentracing.StartSpan("Start")

	// Put root span in context so it will be used in our calls to the client.
	ctx := opentracing.ContextWithSpan(context.Background(), span)
	time.Sleep(time.Duration(20) * time.Millisecond)
	req := &pb.HelloRequest{Name: name}
	res, err := client.SayHello(ctx, req)
	fmt.Println("req: ", req, ", res: ", res)



	span.Finish()

}
