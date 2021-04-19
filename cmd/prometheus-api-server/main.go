package main

import (
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/connect"
	"github.com/powerslider/prometheus-grpc-exporter/pkg/api/server"
	"github.com/powerslider/prometheus-grpc-exporter/pkg/storage"
	"github.com/powerslider/prometheus-grpc-exporter/pkg/transport"
	grpctransport "github.com/powerslider/prometheus-grpc-exporter/pkg/transport/grpc"
	"github.com/powerslider/prometheus-grpc-exporter/pkg/transport/tcp"
	pb "github.com/powerslider/prometheus-grpc-exporter/proto"
	"google.golang.org/grpc"

	cli "github.com/jawher/mow.cli"
)

const (
	appName = "prometheus-api-server"
)

func main() {
	app := cli.App(appName, "")

	grpcPort := app.String(cli.StringOpt{
		Name:   "grpc-port",
		Value:  "8090",
		Desc:   "gRPC Port to listen on",
		EnvVar: "APP_GRPC_PORT",
	})

	//httpPort := app.String(cli.StringOpt{
	//	Name:   "http-port",
	//	Value:  "8081",
	//	Desc:   "HTTP Port to listen on",
	//	EnvVar: "APP_HTTP_PORT",
	//})

	tcpPort := app.String(cli.StringOpt{
		Name:   "tcp-port",
		Value:  "8070",
		Desc:   "TCP Port to listen on",
		EnvVar: "APP_TCP_PORT",
	})

	dataDir := app.String(cli.StringOpt{
		Name:   "data-dir",
		Value:  "/tmp/metrics_store",
		Desc:   "Storage directory for metrics",
		EnvVar: "APP_DATA_DIR",
	})

	// Create a Consul API client
	consulClient, _ := api.NewClient(api.DefaultConfig())

	consulService, _ := connect.NewService(appName, consulClient)
	defer consulService.Close()

	db, err := storage.NewPersistence(*dataDir)
	if err != nil {
		panic(err)
	}
	apiServer := server.NewAPIServer(db)

	grpcTCPListener, err := tcp.NewTCPListener(*grpcPort, consulService.ServerTLSConfig())
	if err != nil {
		panic(err)
	}
	grpcServer := grpctransport.NewGRPCServer(appName, *grpcPort, grpcTCPListener)
	grpcServer.Start(func(s *grpc.Server) {
		pb.RegisterPrometheusServiceServer(s, apiServer)
	})
	//
	//httpServer := httptransport.StartHTTPServer(*httpPort,
	//	httptransport.NewHealthCheckHandler(),
	//)

	tcpListener, err := tcp.NewTCPListener(*tcpPort, consulService.ServerTLSConfig())
	if err != nil {
		panic(err)
	}
	tcpServer := tcp.NewTCPServer(appName, *tcpPort, tcpListener)
	tcpServer.Accept(apiServer.ProcessMetrics)

	transport.WaitForShutdownSignal()

	grpcServer.Shutdown()
	//httptransport.ShutdownHTTPServer(appName, httpServer)
}
