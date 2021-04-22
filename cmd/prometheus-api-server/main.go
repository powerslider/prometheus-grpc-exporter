package main

import (
	"fmt"

	"github.com/powerslider/prometheus-grpc-exporter/pkg/sd"

	"github.com/powerslider/prometheus-grpc-exporter/pkg/api/server"
	"github.com/powerslider/prometheus-grpc-exporter/pkg/storage"
	"github.com/powerslider/prometheus-grpc-exporter/pkg/transport"
	grpctransport "github.com/powerslider/prometheus-grpc-exporter/pkg/transport/grpc"
	httptransport "github.com/powerslider/prometheus-grpc-exporter/pkg/transport/http"
	"github.com/powerslider/prometheus-grpc-exporter/pkg/transport/tcp"
	pb "github.com/powerslider/prometheus-grpc-exporter/proto"
	"google.golang.org/grpc"

	cli "github.com/jawher/mow.cli"
)

const (
	defaultAppName = "prometheus-api-server"
	defaultDataDir = "/tmp/metrics_store"
)

func main() {
	app := cli.App(defaultAppName, "")

	appServiceName := app.String(cli.StringOpt{
		Name:   "app-name",
		Value:  defaultAppName,
		Desc:   "Application Instance Name",
		EnvVar: "APP_SERVICE_NAME",
	})

	appInstanceQualifier := app.String(cli.StringOpt{
		Name:   "app-qualifier",
		Value:  "0",
		Desc:   "Application Instance Qualifier",
		EnvVar: "APP_INSTANCE_QUALIFIER",
	})
	grpcPort := app.String(cli.StringOpt{
		Name:   "grpc-port",
		Value:  "8090",
		Desc:   "gRPC Port to listen on",
		EnvVar: "APP_GRPC_PORT",
	})

	httpPort := app.String(cli.StringOpt{
		Name:   "http-port",
		Value:  "8080",
		Desc:   "HTTP Port to listen on",
		EnvVar: "APP_HTTP_PORT",
	})

	tcpPort := app.String(cli.StringOpt{
		Name:   "tcp-port",
		Value:  "8070",
		Desc:   "TCP Port to listen on",
		EnvVar: "APP_TCP_PORT",
	})

	dataDir := app.String(cli.StringOpt{
		Name:   "data-dir",
		Value:  defaultDataDir,
		Desc:   "Storage directory for scraped Prometheus metrics",
		EnvVar: "APP_INSTANCE_DATA_DIR",
	})

	appInstanceName := fmt.Sprintf("%s-%s", *appServiceName, *appInstanceQualifier)
	httpHealthCheckAddr := fmt.Sprintf("%s:%s", appInstanceName, *httpPort)
	httpServer := httptransport.NewHTTPServer(appInstanceName, *httpPort,
		httptransport.NewHealthCheckHandler(),
	)
	httpServer.Start()

	consulService, err := sd.NewConsulRegistration(*appServiceName, appInstanceName, httpHealthCheckAddr)
	if err != nil {
		panic(err)
	}

	db, err := storage.NewPersistence(*dataDir)
	if err != nil {
		panic(err)
	}
	apiServer := server.NewAPIServer(db)

	grpcTCPListener, err := tcp.NewTCPListener(*grpcPort, consulService.Service.ServerTLSConfig())
	if err != nil {
		panic(err)
	}
	grpcServer := grpctransport.NewGRPCServer(appInstanceName, *grpcPort, grpcTCPListener, func(s *grpc.Server) {
		pb.RegisterPrometheusServiceServer(s, apiServer)
	})
	grpcServer.Start()

	tcpListener, err := tcp.NewTCPListener(*tcpPort, consulService.Service.ServerTLSConfig())
	if err != nil {
		panic(err)
	}
	tcpServer := tcp.NewTCPServer(appInstanceName, *tcpPort, tcpListener)
	tcpServer.Accept(apiServer.ProcessMetrics)

	transport.WaitForShutdownSignal()

	grpcServer.Shutdown()
	httpServer.Shutdown()
}
