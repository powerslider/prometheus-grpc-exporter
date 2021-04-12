package prometheus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"time"

	pb "github.com/powerslider/prometheus-grpc-exporter/proto"

	"github.com/golang/snappy"
	//nolint:staticcheck
	"github.com/golang/protobuf/proto"
	"github.com/prometheus/prometheus/prompb"
)

var CurrentMetrics []*pb.TimeSeries

type MetricsStatus struct {
	Msg          string    `json:"message"`
	LastModified time.Time `json:"last_modified"`
}

type Options map[string]string

type MetricsHandler struct {
	Options Options
}

func (mh *MetricsHandler) RemoteWriteHandler(w http.ResponseWriter, r *http.Request) {
	req, err := decodeWriteRequest(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var metrics []prompb.TimeSeries
	for _, ts := range req.Timeseries {
		var samples []prompb.Sample
		for _, s := range ts.Samples {
			if math.IsNaN(s.Value) {
				continue
			}
			samples = append(samples, s)
		}
		metrics = append(metrics, prompb.TimeSeries{Labels: ts.Labels, Samples: samples})
	}

	payload, err := json.Marshal(metrics)
	log.Println(string(payload))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Fatal("error serializing metrics response: ", err)
		return
	}

	resp, err := sendLBRequest(payload, mh.Options)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Fatal("error sending metrics to LB: ", err)
		return
	}
	defer resp.Body.Close()

	respondWithMetricsStatus("New metrics consumed", w)
}

func sendLBRequest(payload []byte, options Options) (*http.Response, error) {
	url := fmt.Sprintf("http://%s/metrics", options["lb_host"])
	resp, err := http.Post(url, "application/json", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (mh *MetricsHandler) ConsumeMetricsHandler(w http.ResponseWriter, r *http.Request) {
	reqBodyBytes, _ := requestBodyToBytes(r)
	err := json.Unmarshal(reqBodyBytes, &CurrentMetrics)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Fatal("error deserializing metrics response: ", err)
		return
	}
	log.Println(CurrentMetrics)
	respondWithMetricsStatus("New metrics processed", w)
}

func requestBodyToBytes(r *http.Request) ([]byte, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// DecodeWriteRequest from an io.Reader into a prompb.WriteRequest, handling
// snappy decompression.
func decodeWriteRequest(r io.Reader) (*prompb.WriteRequest, error) {
	compressed, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	reqBuf, err := snappy.Decode(nil, compressed)
	if err != nil {
		return nil, err
	}

	var req prompb.WriteRequest
	if err := proto.Unmarshal(reqBuf, &req); err != nil {
		return nil, err
	}

	return &req, nil
}

func respondWithMetricsStatus(statusMessage string, w http.ResponseWriter) {
	currentStatus := MetricsStatus{Msg: statusMessage, LastModified: time.Now()}
	currentStatusJSON, err := json.Marshal(currentStatus)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal("error serializing metrics status: ", err)
		return
	}
	if _, err := w.Write(currentStatusJSON); err != nil {
		w.WriteHeader(http.StatusCreated)
		log.Fatal("metrics server error: ", err)
	}
}
