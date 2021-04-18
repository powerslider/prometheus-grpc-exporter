package prometheus

import (
	"io"
	"io/ioutil"
	"math"

	"github.com/golang/snappy"

	pb "github.com/powerslider/prometheus-grpc-exporter/proto"

	"github.com/golang/protobuf/proto"
	"github.com/prometheus/prometheus/prompb"
)

func ProcessMetrics(r io.Reader) (*pb.MetricsBatch, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	req, err := deserializeWriteRequest(data)
	if err != nil {
		return nil, err
	}
	metricsBatch := deserializeMetrics(req.Timeseries)
	return metricsBatch, nil
}

// DecodeWriteRequest from an io.Reader into a prompb.WriteRequest, handling
// snappy decompression.
func DecodeWriteRequest(r io.Reader) (*prompb.WriteRequest, error) {
	compressed, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	reqBuf, err := snappy.Decode(nil, compressed)
	if err != nil {
		return nil, err
	}

	return deserializeWriteRequest(reqBuf)
}

func deserializeWriteRequest(reqBuf []byte) (*prompb.WriteRequest, error) {
	var req prompb.WriteRequest
	if err := proto.Unmarshal(reqBuf, &req); err != nil {
		return nil, err
	}
	return &req, nil
}
func deserializeMetrics(inputMetrics []*prompb.TimeSeries) *pb.MetricsBatch {
	metrics := make([]*pb.TimeSeries, 0)
	for _, ts := range inputMetrics {
		labels := make([]*pb.Label, 0)
		for _, l := range ts.Labels {
			if l.Name == "" || l.Value == "" {
				continue
			}
			labels = append(labels, &pb.Label{
				Name:  l.Name,
				Value: l.Value,
			})
		}
		samples := make([]*pb.Sample, 0)
		for _, s := range ts.Samples {
			if math.IsNaN(s.Value) {
				continue
			}
			samples = append(samples, &pb.Sample{
				Value:     s.Value,
				Timestamp: s.Timestamp,
			})
		}
		if len(labels) == 0 || len(samples) == 0 {
			continue
		}
		metrics = append(metrics, &pb.TimeSeries{Labels: labels, Samples: samples})
	}
	return &pb.MetricsBatch{
		TimeSeries: metrics,
	}
}
