package main

import (
	"fmt"
	sd "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v2"
	"google.golang.org/grpc"

	"log"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9000))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	sd.RegisterSecretDiscoveryServiceServer(grpcServer, &MySDS{})
	grpcServer.Serve(lis)
}
