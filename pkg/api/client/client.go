package client

import (
	"encoding/json"
	"io"
	"log"

	pb "github.com/powerslider/prometheus-grpc-exporter/proto"
)

func GetMetricsResponse(stream pb.PrometheusService_GetMetricsClient) {
	done := make(chan bool)

	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				done <- true
				return
			}
			if err != nil {
				log.Fatalf("cannot receive %v", err)
			}

			currentMetric, err := json.Marshal(resp)
			if err != nil {
				log.Fatalf("cannot serialize time series metric to json: %v", err)
			}
			log.Printf("received metric: %s", currentMetric)
		}
	}()

	<-done
	log.Printf("finished")
}
