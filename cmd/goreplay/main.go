package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/powerslider/prometheus-grpc-exporter/pkg/goreplay"

	cli "github.com/jawher/mow.cli"
)

func main() {
	app := cli.App("goreplay", "")

	port := app.String(cli.StringOpt{
		Name:   "port",
		Value:  "8082",
		Desc:   "Port to listen on",
		EnvVar: "APP_PORT",
	})

	http.HandleFunc("/receive", goreplay.PrometheusRemoteWriteHandler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", *port), nil))
}
