package main

import (
	"github.com/powerslider/prometheus-grpc-exporter/pkg/api/server"
	"github.com/powerslider/prometheus-grpc-exporter/pkg/prometheus"
	"github.com/powerslider/prometheus-grpc-exporter/pkg/transport"
	grpctransport "github.com/powerslider/prometheus-grpc-exporter/pkg/transport/grpc"
	httptransport "github.com/powerslider/prometheus-grpc-exporter/pkg/transport/http"
	pb "github.com/powerslider/prometheus-grpc-exporter/proto"
	"google.golang.org/grpc"

	cli "github.com/jawher/mow.cli"
)

const appName = "prometheus-grpc-server"

func main() {
	app := cli.App(appName, "")

	grpcPort := app.String(cli.StringOpt{
		Name:   "grpc-port",
		Value:  "8090",
		Desc:   "Port to listen on",
		EnvVar: "APP_GRPC_PORT",
	})

	httpPort := app.String(cli.StringOpt{
		Name:   "http-port",
		Value:  "8080",
		Desc:   "Port to listen on",
		EnvVar: "APP_HTTP_PORT",
	})

	httpServer := httptransport.StartHTTPServer(*httpPort,
		httptransport.Handler{
			Path:        "/metrics",
			HandlerFunc: prometheus.ConsumeMetricsHandler,
		},
		httptransport.NewHealthCheckHandler(),
	)

	grpcServer := grpctransport.StartGRPCServer(*grpcPort, func(s *grpc.Server) {
		pb.RegisterPrometheusServiceServer(s, server.Server{})
	})

	transport.WaitForShutdownSignal()

	httptransport.ShutdownHTTPServer(appName, httpServer)
	grpctransport.ShutdownGRPCServer(appName, grpcServer)
}
