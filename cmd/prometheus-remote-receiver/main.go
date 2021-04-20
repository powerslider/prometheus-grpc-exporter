package main

import (
	"net"

	cli "github.com/jawher/mow.cli"
	"github.com/powerslider/prometheus-grpc-exporter/pkg/prometheus"
	"github.com/powerslider/prometheus-grpc-exporter/pkg/transport"
	httptransport "github.com/powerslider/prometheus-grpc-exporter/pkg/transport/http"
	"github.com/powerslider/prometheus-grpc-exporter/pkg/transport/tcp"
)

const (
	appName = "prometheus-remote-receiver"
)

func main() {
	app := cli.App(appName, "")

	httpPort := app.String(cli.StringOpt{
		Name:   "http-port",
		Value:  "8080",
		Desc:   "Port to listen on",
		EnvVar: "APP_HTTP_PORT",
	})

	metricsForwarderAddr := app.String(cli.StringOpt{
		Name:   "metrics-forwarder-addr",
		Value:  "metrics-forwarder:9999",
		Desc:   "Metrics Forwarder Address",
		EnvVar: "METRICS_FORWARDER_ADDR",
	})

	mh := prometheus.NewMetricsHandler(
		prometheus.Options{
			"tcp_forwarder_addr": *metricsForwarderAddr,
		},
		prometheus.MetricsStore{},
	)

	httpServer := httptransport.NewHTTPServer(appName, *httpPort,
		httptransport.Handler{
			Path:        "/",
			HandlerFunc: mh.RemoteWriteHandler,
		},
		httptransport.NewHealthCheckHandler(),
	)
	httpServer.Start()

	tcpClient := tcp.NewTCPClient(*metricsForwarderAddr)

	tcpClient.Connect(func(conn net.Conn) {
		for {
			if len(mh.Store) > 0 {
				conn.Write(mh.Store)
				mh.Store = prometheus.MetricsStore{}
			}
		}
	})

	transport.WaitForShutdownSignal()

	httpServer.Shutdown()
}
