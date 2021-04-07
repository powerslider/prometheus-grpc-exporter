package server

import (
	"fmt"
	"log"
	"sync"
	"time"

	pb "github.com/powerslider/prometheus-grpc-exporter/proto"
)

type Server struct{}

func (s Server) GetMetrics(in *pb.GetMetricsRequest, srv pb.PrometheusService_GetMetricsServer) error {
	log.Printf("consume metrics for id : %d", in.Id)

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(count int64) {
			defer wg.Done()
			time.Sleep(time.Duration(count) * time.Second)
			resp := pb.MetricsResponse{Result: fmt.Sprintf("Request #%d For Id:%d", count, in.Id)}
			if err := srv.Send(&resp); err != nil {
				log.Printf("send error %v", err)
			}
			log.Printf("finishing request number : %d", count)
		}(int64(i))
	}

	wg.Wait()
	return nil
}
