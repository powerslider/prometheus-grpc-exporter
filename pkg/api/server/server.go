package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"sync"
	"time"

	"github.com/powerslider/prometheus-grpc-exporter/pkg/prometheus"
	pb "github.com/powerslider/prometheus-grpc-exporter/proto"
	"google.golang.org/protobuf/proto"
)

type Server struct {
}

func (s *Server) GetMetrics(in *pb.GetMetricsRequest, srv pb.PrometheusService_GetMetricsServer) error {
	log.Printf("consume metrics for id : %d", in.Id)

	var metricsBatch pb.MetricsBatch
	err := s.readMetrics("./storage/1618703265.bin", &metricsBatch)
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
	if err := s.storeMetrics(metricsBatch); err != nil {
		log.Fatal("error storing metrics batch:", err)
	}
}

func (s *Server) storeMetrics(message proto.Message) error {
	data, err := json.Marshal(message)
	//data, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("cannot marshal proto message to binary: %w", err)
	}

	fileName := fmt.Sprintf("./storage/%d.bin", time.Now().Unix())
	err = ioutil.WriteFile(fileName, data, 0644)
	if err != nil {
		return fmt.Errorf("cannot write binary data to file: %w", err)
	}

	return nil
}

func (s *Server) readMetrics(filename string, message proto.Message) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("cannot read binary data from file: %w", err)
	}

	err = proto.Unmarshal(data, message)
	if err != nil {
		return fmt.Errorf("cannot unmarshal binary to proto message: %w", err)
	}

	return nil
}
