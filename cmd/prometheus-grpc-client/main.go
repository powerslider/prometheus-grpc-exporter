package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	cli "github.com/jawher/mow.cli"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"github.com/powerslider/prometheus-grpc-exporter/pkg/api/client"
	pb "github.com/powerslider/prometheus-grpc-exporter/proto"
	"google.golang.org/grpc"
)

const appName = "prometheus-grpc-client"

func main() {
	app := cli.App(appName, "")

	serverHost := app.String(cli.StringOpt{
		Name:   "server-host",
		Value:  "localhost:8090",
		Desc:   "Server Host Address",
		EnvVar: "SERVER_HOST",
	})

	subscribedMetric := app.String(cli.StringOpt{
		Name:   "subscribed-metric",
		Value:  "",
		Desc:   "Specifically received metric",
		EnvVar: "SUBSCRIBED_METRIC",
	})

	rand.Seed(time.Now().Unix())

	// dial server
	conn, err := grpc.Dial(*serverHost, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("can not connect with server %v", err)
	}

	// create stream
	c := pb.NewPrometheusServiceClient(conn)
	in := &pb.GetMetricsRequest{Id: 1}
	stream, err := c.GetMetrics(context.Background(), in)
	if err != nil {
		log.Fatalf("open stream error %v", err)
	}

	client.GetMetricsResponse(stream, subscribedMetric)
}
