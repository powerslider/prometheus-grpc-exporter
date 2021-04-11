package client

import (
	"encoding/json"
	"io"
	"log"

	pb "github.com/powerslider/prometheus-grpc-exporter/proto"
)

const metricLabelName = "__name__"

func GetMetricsResponse(stream pb.PrometheusService_GetMetricsClient, subscribedMetric *string) {
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

			if len(*subscribedMetric) > 0 {
				for _, l := range resp.Labels {
					if l.Name == metricLabelName && l.Value == *subscribedMetric {
						outputMetric(resp)
					}
				}
			} else {
				outputMetric(resp)
			}
		}
	}()

	<-done
	log.Printf("finished")
}

func outputMetric(resp *pb.TimeSeries) {
	currentMetric, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("cannot serialize time series metric to json: %v", err)
	}
	log.Printf("received metric: %s\n", currentMetric)
}
