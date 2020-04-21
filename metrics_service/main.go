package main

import (
	"flag"
	"fmt"
	v2 "github.com/envoyproxy/go-control-plane/envoy/service/metrics/v2"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 10001))
	log.Println("metrics server listen to 10001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	v2.RegisterMetricsServiceServer(grpcServer, &MyMetricsServer{})
	_ = grpcServer.Serve(lis)
}
