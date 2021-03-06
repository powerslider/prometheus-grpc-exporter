package server

import (
	"log"
	"net"
	"sync"
	"time"

	"github.com/powerslider/prometheus-grpc-exporter/pkg/storage"

	"github.com/powerslider/prometheus-grpc-exporter/pkg/prometheus"
	pb "github.com/powerslider/prometheus-grpc-exporter/proto"
)

type Server struct {
	Storage *storage.Persistence
}

func NewAPIServer(storage *storage.Persistence) *Server {
	return &Server{
		Storage: storage,
	}
}

func (s *Server) GetMetrics(in *pb.GetMetricsRequest, srv pb.PrometheusService_GetMetricsServer) error {
	log.Printf("consume metrics for id : %d", in.Id)

	var metricsBatch pb.MetricsBatch
	err := s.Storage.Read(&metricsBatch)
	if err != nil {
		log.Fatal("error reading metrics batch: ", err)
	}
	var wg sync.WaitGroup
	for i, m := range metricsBatch.TimeSeries {
		wg.Add(1)
		go func(count int64, currentMetric *pb.TimeSeries) {
			defer wg.Done()
			time.Sleep(time.Duration(count) * time.Second)
			if err := srv.Send(currentMetric); err != nil {
				log.Printf("send error %v", err)
			}
			log.Printf("finishing request number : %d", count)
		}(int64(i), m)
	}

	wg.Wait()
	return nil
}

func (s *Server) ProcessMetrics(conn net.Conn) {
	metricsBatch, err := prometheus.ProcessMetrics(conn)
	if err != nil {
		log.Fatal("error processing incoming metrics: ", err)
		return
	}
	if err := s.Storage.Save(metricsBatch); err != nil {
		log.Fatal("error storing metrics batch:", err)
	}
}
