package main

import (
	"fmt"
	pb "gsf/src/QosBreaker/helloworld"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gsf/src/QosRateLimit/rateLimit"
	"log"
	"net"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	fmt.Println("######### get client request name :"+in.Name)
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	/*
	interceptor:
	grpc服务端提供interceptor功能，可以在服务端接收到请求时优先对请求中的数据做一些处理后，再转交给指定的服务处理并响应，
	功能类似middleware，适合做处理验证，日志等流程。

	在自定义Token认证的示例中，认证信息是由每个服务中的方法处理并认证的，如果有大量的接口方法，这种姿势就太蛋疼了，
	每个接口实现都要先处理认证信息。这个时候interceptor就站出来解决了这个问题，可以在请求被转到具体接口之前处理认证信息。

	go-grpc-middleware:
	这个项目对grpc的interceptor进行了封装，支持多个拦截器的链式组装。
	 */


	/*
	限流
	ratelimit interceptor
	*/
	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			qosRateLimit.LimitRateServerInterceptor(),
			)))





	//s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}



}

