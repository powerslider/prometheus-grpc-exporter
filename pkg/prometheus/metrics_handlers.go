package prometheus

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/golang/snappy"
	//nolint:staticcheck
	"github.com/golang/protobuf/proto"
	"github.com/prometheus/prometheus/prompb"
)

var CurrentMetrics []prompb.TimeSeries

type MetricsStatus struct {
	Msg          string    `json:"message"`
	LastModified time.Time `json:"last_modified"`
}

func RemoteWriteHandler(w http.ResponseWriter, r *http.Request) {
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

	resp, err := json.Marshal(metrics)
	log.Println(string(resp))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Fatal("error serializing metrics response: ", err)
		return
	}
	//TODO: call LB

	respondWithMetricsStatus("New metrics consumed", w)
}

func ConsumeMetricsHandler(w http.ResponseWriter, r *http.Request) {
	var metrics []prompb.TimeSeries
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	reqBodyBytes := buf.Bytes()
	log.Println(buf.String())
	err := json.Unmarshal(reqBodyBytes, &CurrentMetrics)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Fatal("error deserializing metrics response: ", err)
		return
	}
	log.Println(metrics)
	respondWithMetricsStatus("New metrics processed", w)
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
	currentStatusJson, err := json.Marshal(currentStatus)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal("error serializing metrics status: ", err)
		return
	}
	if _, err := w.Write(currentStatusJson); err != nil {
		w.WriteHeader(http.StatusCreated)
		log.Fatal("metrics server error: ", err)
	}
}
