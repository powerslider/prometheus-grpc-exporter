package main

import (
	"fmt"
	"strings"

	cli "github.com/jawher/mow.cli"
	"github.com/powerslider/prometheus-grpc-exporter/pkg/transport"
	httptransport "github.com/powerslider/prometheus-grpc-exporter/pkg/transport/http"
	"github.com/powerslider/prometheus-grpc-exporter/pkg/transport/tcp"
)

const (
	appName = "metrics-forwarder"
)

func main() {
	app := cli.App(appName, "")

	localPort := app.String(cli.StringOpt{
		Name:   "local-addr",
		Value:  "9999",
		Desc:   "Local Address",
		EnvVar: "FORWARD_LOCAL_PORT",
	})

	httpPort := app.String(cli.StringOpt{
		Name:   "http-port",
		Value:  "8080",
		Desc:   "HTTP Port to listen on",
		EnvVar: "APP_HTTP_PORT",
	})

	remoteAddresses := app.String(cli.StringOpt{
		Name:   "forward-remote-addresses",
		Value:  "prometheus-api-server-1;prometheus-api-server-2",
		Desc:   "Remote Address",
		EnvVar: "FORWARD_REMOTE_ADDRESSES",
	})

	httpServer := httptransport.NewHTTPServer(appName, *httpPort,
		httptransport.NewHealthCheckHandler(),
	)
	httpServer.Start()

	addrs := strings.Split(*remoteAddresses, ";")
	httpHealthCheckAddr := fmt.Sprintf("%s:%s", appName, *httpPort)
	forwarder := tcp.NewTCPForwarder(appName, *localPort, addrs, httpHealthCheckAddr)
	forwarder.Accept()

	transport.WaitForShutdownSignal()

	httpServer.Shutdown()
}
