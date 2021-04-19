package main

import (
	cli "github.com/jawher/mow.cli"
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

	remoteAddr := app.String(cli.StringOpt{
		Name:   "remote-addr",
		Value:  "prometheus-api-server",
		Desc:   "Remote Address",
		EnvVar: "FORWARD_REMOTE_ADDR",
	})

	forwarder := tcp.NewTCPForwarder(appName, *localPort, *remoteAddr)
	forwarder.Accept()
}
