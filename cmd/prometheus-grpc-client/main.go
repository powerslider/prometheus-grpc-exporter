package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc/resolver"

	cli "github.com/jawher/mow.cli"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"github.com/powerslider/prometheus-grpc-exporter/pkg/api/client"
	pb "github.com/powerslider/prometheus-grpc-exporter/proto"
	"github.com/simplesurance/grpcconsulresolver/consul"
	"google.golang.org/grpc"
)

const appName = "prometheus-grpc-client"

func init() {
	// Register the consul consul at the grpc-go library
	resolver.Register(consul.NewBuilder())
}

func main() {
	app := cli.App(appName, "")

	consulHost := app.String(cli.StringOpt{
		Name:   "consul-host",
		Value:  "consul-server:8500",
		Desc:   "Consul Host Address",
		EnvVar: "CONSUL_HOST",
	})

	subscribedMetric := app.String(cli.StringOpt{
		Name:   "subscribed-metric",
		Value:  "",
		Desc:   "Specifically received metric",
		EnvVar: "SUBSCRIBED_METRIC",
	})

	// dial server
	connStr := fmt.Sprintf("consul://%s/prometheus-api-server-1", *consulHost)
	conn, err := grpc.Dial(connStr, grpc.WithInsecure(), grpc.WithBalancerName("round_robin"))
	//conn, err := grpc.Dial("consul://metrics", grpc.WithInsecure())
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

	client.GetMetricsResponse(stream, *subscribedMetric)
}
