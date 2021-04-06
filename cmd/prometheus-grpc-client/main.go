package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	cli "github.com/jawher/mow.cli"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"github.com/powerslider/prometheus-grpc-exporter/pkg/client"
	pb "github.com/powerslider/prometheus-grpc-exporter/proto"
	"google.golang.org/grpc"
)

func main() {
	app := cli.App("prometheus-grpc-client", "")

	serverHost := app.String(cli.StringOpt{
		Name:   "server-host",
		Value:  "localhost:8080",
		Desc:   "Server Host Address",
		EnvVar: "SERVER_HOST",
	})
	rand.Seed(time.Now().Unix())

	// dial server
	conn, err := grpc.Dial(*serverHost, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("can not connect with server %v", err)
	}

	// create stream
	c := pb.NewPrometheusServiceClient(conn)
	in := &pb.ConsumeMetricsRequest{Id: 1}
	stream, err := c.ConsumeMetrics(context.Background(), in)
	if err != nil {
		log.Fatalf("open stream error %v", err)
	}

	client.GetMetricsResponse(stream)
}
