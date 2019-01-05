package main

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "gsf/src/QosTracer/helloworld"
	"log"
	"net"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"

	"github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	fmt.Println("######### get client request name :" + in.Name)
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {



	//zipkin
	collector, err := zipkin.NewHTTPCollector("http://localhost:9411/api/v1/spans")
	if err != nil {
		panic(err)
		return
	}

	tracer, err := zipkin.NewTracer(
		zipkin.NewRecorder(collector, false, "localhost:0", "grpc_server"),
		zipkin.ClientServerSameSpan(true),
		zipkin.TraceID128Bit(true),
	)


	if err != nil {
		panic(err)
	}
	opentracing.InitGlobalTracer(tracer)
	l, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(tracer, otgrpc.LogPayloads())),
	)

	//s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	s.Serve(l)

}
