package main

import (
	"github.com/powerslider/prometheus-grpc-exporter/pkg/prometheus"
	"github.com/powerslider/prometheus-grpc-exporter/pkg/transport"
	httptransport "github.com/powerslider/prometheus-grpc-exporter/pkg/transport/http"

	cli "github.com/jawher/mow.cli"
)

const appName = "prometheus-remote-writer"

func main() {
	app := cli.App(appName, "")

	httpPort := app.String(cli.StringOpt{
		Name:   "http-port",
		Value:  "8080",
		Desc:   "Port to listen on",
		EnvVar: "APP_HTTP_PORT",
	})

	httpServer := httptransport.StartHTTPServer(*httpPort,
		httptransport.Handler{
			Path:        "/",
			HandlerFunc: prometheus.RemoteWriteHandler,
		},
	)

	transport.WaitForShutdownSignal()
	httptransport.ShutdownHTTPServer(appName, httpServer)
}
