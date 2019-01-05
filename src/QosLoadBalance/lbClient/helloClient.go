package main

import (
	"flag"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "gsf/src/QosLoadBalance/helloworld"
	"gsf/src/QosLoadBalance/lbMdw"
	"strconv"
	"time"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
	serviceName = "qosCircuitBreakerService"
)

var (
	serv = flag.String("service", "hello_service", "service name")
	reg = flag.String("reg", "http://127.0.0.1:2379", "register etcd address")
)


/*
进程内LB
 */

func main() {

	flag.Parse()
	r := grpclb.NewResolver(*serv)
	//grpc自己的balancer
	b := grpc.RoundRobin(r)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	conn, err := grpc.DialContext(ctx, *reg, grpc.WithInsecure(), grpc.WithBalancer(b))
	if err != nil {
		panic(err)
	}



	ticker := time.NewTicker(1 * time.Second)
	for t := range ticker.C {
		client := pb.NewGreeterClient(conn)
		resp, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "world " + strconv.Itoa(t.Second())})
		if err == nil {
			fmt.Printf("%v: Reply is %s\n", t, resp.Message)
		}
	}



}
