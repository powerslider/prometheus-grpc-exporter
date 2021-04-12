package main

import (
	"github.com/powerslider/prometheus-grpc-exporter/pkg/prometheus"
	"github.com/powerslider/prometheus-grpc-exporter/pkg/transport"
	httptransport "github.com/powerslider/prometheus-grpc-exporter/pkg/transport/http"

	cli "github.com/jawher/mow.cli"
)

const appName = "prometheus-remote-receiver"

func main() {
	app := cli.App(appName, "")

	httpPort := app.String(cli.StringOpt{
		Name:   "http-port",
		Value:  "8080",
		Desc:   "Port to listen on",
		EnvVar: "APP_HTTP_PORT",
	})

	lbHost := app.String(cli.StringOpt{
		Name:   "lb-host",
		Value:  "fabio:9999",
		Desc:   "Load Balancer Host",
		EnvVar: "LB_HOST",
	})

	mh := prometheus.MetricsHandler{
		Options: map[string]string{
			"lb_host": *lbHost,
		},
	}
	httpServer := httptransport.StartHTTPServer(*httpPort,
		httptransport.Handler{
			Path:        "/",
			HandlerFunc: mh.RemoteWriteHandler,
		},
		httptransport.NewHealthCheckHandler(),
	)

	transport.WaitForShutdownSignal()
	httptransport.ShutdownHTTPServer(appName, httpServer)
}
