package http

import (
	"log"
	"net/http"
)

const healthEndpoint = "/__health"

func NewHealthCheckHandler() Handler {
	return Handler{
		Path:        healthEndpoint,
		HandlerFunc: checkHealthHandler,
	}
}
func checkHealthHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("OK")); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal("server error: ", err)
	}
}
