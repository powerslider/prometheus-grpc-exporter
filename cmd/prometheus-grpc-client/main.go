package main

import (
	"context"
	"log"

	"google.golang.org/grpc"

	cli "github.com/jawher/mow.cli"
	"github.com/powerslider/prometheus-grpc-exporter/pkg/api/client"
	pb "github.com/powerslider/prometheus-grpc-exporter/proto"
)

const appName = "prometheus-grpc-client"

func main() {
	app := cli.App(appName, "")

	subscribedMetric := app.String(cli.StringOpt{
		Name:   "subscribed-metric",
		Value:  "",
		Desc:   "Specifically received metric",
		EnvVar: "SUBSCRIBED_METRIC",
	})

	lbAddr := app.String(cli.StringOpt{
		Name:   "lb-addr",
		Value:  "haproxy:8090",
		Desc:   "Haproxy gRPC LB Address",
		EnvVar: "LB_ADDR",
	})

	//rr := grpc.RoundRobin(grpcsrvlb.New(srv.NewGoResolver(2 * time.Second)))
	//conn, err := grpc.Dial("grpc.my_service.my_cluster.internal.example.com", grpc.WithBalancer(rr))
	// dial server
	conn, err := grpc.Dial(*lbAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("cannot connect to server %v", err)
	}

	// create stream
	c := pb.NewPrometheusServiceClient(conn)
	in := &pb.GetMetricsRequest{Id: 1}
	stream, err := c.GetMetrics(context.Background(), in)
	if err != nil {
		log.Fatalf("open stream error %v", err)
	}

	client.GetMetricsResponse(stream, *subscribedMetric)
}
