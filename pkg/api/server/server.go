package server

import (
	"log"
	"sync"
	"time"

	"github.com/powerslider/prometheus-grpc-exporter/pkg/prometheus"

	pb "github.com/powerslider/prometheus-grpc-exporter/proto"
)

type Server struct{}

func (s Server) GetMetrics(in *pb.GetMetricsRequest, srv pb.PrometheusService_GetMetricsServer) error {
	log.Printf("consume metrics for id : %d", in.Id)

	// NOTE: use this only for testing purposes
	//var currentMetrics []*pb.TimeSeries
	//data, _ := ioutil.ReadFile("testdata/test_metrics.json")
	//_ = json.Unmarshal(data, &currentMetrics)
	var wg sync.WaitGroup
	for i, m := range prometheus.CurrentMetrics {
		//for i, m := range currentMetrics {
		wg.Add(1)
		go func(currentMetric *pb.TimeSeries, count int64) {
			defer wg.Done()
			time.Sleep(time.Duration(count) * time.Second)
			if err := srv.Send(currentMetric); err != nil {
				log.Printf("send error %v", err)
			}
			log.Printf("finishing request number : %d", count)
		}(m, int64(i))
	}

	wg.Wait()
	return nil
}
