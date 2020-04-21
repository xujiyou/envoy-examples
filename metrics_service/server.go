package main

import (
	v2 "github.com/envoyproxy/go-control-plane/envoy/service/metrics/v2"
	"io"
	"log"
)

type MyMetricsServer struct{}

func (myServer *MyMetricsServer) StreamMetrics(server v2.MetricsService_StreamMetricsServer) error {
	for {
		message, err := server.Recv()
		if err == io.EOF {
			_ = server.SendAndClose(&v2.StreamMetricsResponse{})
			return nil
		}
		if err != nil {
			log.Printf("failed to recv: %v", err)
			return err
		}
		log.Println(message)
	}
}
