package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/powerslider/prometheus-grpc-exporter/pkg/client"
	pb "github.com/powerslider/prometheus-grpc-exporter/proto"
	"google.golang.org/grpc"
)

func main() {
	rand.Seed(time.Now().Unix())

	// dial server
	conn, err := grpc.Dial("localhost:50005", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("can not connect with server %v", err)
	}

	// create stream
	c := pb.NewPrometheusServiceClient(conn)
	in := &pb.ConsumeMetricsRequest{Id: 1}
	stream, err := c.ConsumeMetrics(context.Background(), in)
	if err != nil {
		log.Fatalf("openn stream error %v", err)
	}

	client.GetMetricsResponse(stream)
}
