package prometheus

import (
	"encoding/json"
	"log"
	"net/http"
)

func ConsumeMetricsHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := json.Marshal(r.Body)
	log.Println(string(resp))
	if err != nil {
		log.Fatal("error serializing metrics response: ", err)
		return
	}

	if _, err := w.Write(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal("metrics server error: ", err)
		return
	}
}
