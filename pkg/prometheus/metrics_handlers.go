package prometheus

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type MetricsStatus struct {
	Msg          string    `json:"message"`
	LastModified time.Time `json:"last_modified"`
}

type Options map[string]string
type MetricsStore []byte
type MetricsHandler struct {
	Options Options
	Store   MetricsStore
}

func NewMetricsHandler(options Options, store []byte) *MetricsHandler {
	return &MetricsHandler{
		Options: options,
		Store:   store,
	}
}
func (mh *MetricsHandler) RemoteWriteHandler(w http.ResponseWriter, r *http.Request) {
	req, err := DecodeWriteRequest(r.Body)
	if err != nil {
		log.Fatal("error processing incoming metrics: ", err)
		return
	}
	data, err := req.Marshal()
	if err != nil {
		log.Fatal("cannot marshal proto message to binary: %w", err)

	}
	mh.Store = data
	respondWithMetricsStatus("New metrics consumed", w)
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
