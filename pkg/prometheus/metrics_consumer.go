package prometheus

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/prometheus/prompb"
)

var CurrentMetrics []prompb.TimeSeries

type MetricsStatus struct {
	Msg          string    `json:"message"`
	LastModified time.Time `json:"last_modified"`
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
	currentStatus := MetricsStatus{Msg: "New metrics added.", LastModified: time.Now()}
	currentStatusJson, err := json.Marshal(currentStatus)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal("error serializing metrics status: ", err)
		return
	}
	if _, err := w.Write(currentStatusJson); err != nil {
		w.WriteHeader(http.StatusCreated)
		log.Fatal("metrics server error: ", err)
		return
	}
}
