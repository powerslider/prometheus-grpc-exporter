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

	tcpForwarderAddr := app.String(cli.StringOpt{
		Name:   "tcp-forwarder-addr",
		Value:  "tcp-forwarder:9999",
		Desc:   "TCP Forwarder Address",
		EnvVar: "TCP_FORWARDER_ADDR",
	})

	mh := prometheus.NewMetricsHandler(
		prometheus.Options{
			"tcp_forwarder_addr": *tcpForwarderAddr,
		},
		prometheus.MetricsStore{},
	)

	httpServer := httptransport.StartHTTPServer(*httpPort,
		httptransport.Handler{
			Path:        "/",
			HandlerFunc: mh.RemoteWriteHandler,
		},
		httptransport.NewHealthCheckHandler(),
	)

	tcpClient := tcp.NewTCPClient(*tcpForwarderAddr)

	tcpClient.Connect(func(conn net.Conn) {
		if len(mh.Store) > 0 {
			conn.Write(mh.Store)
			mh.Store = prometheus.MetricsStore{}
		}
	})

	transport.WaitForShutdownSignal()
	httptransport.ShutdownHTTPServer(appName, httpServer)
}
