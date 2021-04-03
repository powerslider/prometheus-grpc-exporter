package client

import (
	"io"
	"log"

	pb "github.com/powerslider/prometheus-grpc-exporter/proto"
)

func GetMetricsResponse(stream pb.PrometheusService_ConsumeMetricsClient) {
	//ctx := stream.Context()
	done := make(chan bool)

	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				done <- true //close(done)
				return
			}
			if err != nil {
				log.Fatalf("can not receive %v", err)
			}
			log.Printf("Resp received: %s", resp.Result)
		}
	}()

	<-done
	log.Printf("finished")
}
