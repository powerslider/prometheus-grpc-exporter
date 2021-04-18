package main

import (
	"github.com/powerslider/prometheus-grpc-exporter/pkg/api/server"
	"github.com/powerslider/prometheus-grpc-exporter/pkg/transport"
	grpctransport "github.com/powerslider/prometheus-grpc-exporter/pkg/transport/grpc"
	httptransport "github.com/powerslider/prometheus-grpc-exporter/pkg/transport/http"
	"github.com/powerslider/prometheus-grpc-exporter/pkg/transport/tcp"
	pb "github.com/powerslider/prometheus-grpc-exporter/proto"
	"google.golang.org/grpc"

	cli "github.com/jawher/mow.cli"
)

const appName = "prometheus-api-server"

func main() {
	app := cli.App(appName, "")

	grpcPort := app.String(cli.StringOpt{
		Name:   "grpc-port",
		Value:  "8090",
		Desc:   "gRPC Port to listen on",
		EnvVar: "APP_GRPC_PORT",
	})

	httpPort := app.String(cli.StringOpt{
		Name:   "http-port",
		Value:  "8081",
		Desc:   "HTTP Port to listen on",
		EnvVar: "APP_HTTP_PORT",
	})

	tcpPort := app.String(cli.StringOpt{
		Name:   "tcp-port",
		Value:  "8070",
		Desc:   "TCP Port to listen on",
		EnvVar: "APP_TCP_PORT",
	})

	httpServer := httptransport.StartHTTPServer(*httpPort,
		httptransport.NewHealthCheckHandler(),
	)

	apiServer := &server.Server{}
	grpcServer := grpctransport.StartGRPCServer(*grpcPort, func(s *grpc.Server) {
		pb.RegisterPrometheusServiceServer(s, apiServer)
	})

	tcpServer := tcp.NewTCPServer(*tcpPort)
	tcpServer.Accept(apiServer.ProcessMetrics)

	transport.WaitForShutdownSignal()

	httptransport.ShutdownHTTPServer(appName, httpServer)
	grpctransport.ShutdownGRPCServer(appName, grpcServer)
}
