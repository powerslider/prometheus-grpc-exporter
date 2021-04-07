package http

import (
	"context"
	"log"
	"net/http"
	"time"
)

type Handler struct {
	HandlerFunc func(http.ResponseWriter, *http.Request)
	Path        string
}

func StartHTTPServer(port string, handlers ...Handler) *http.Server {
	serveMux := http.NewServeMux()
	for _, h := range handlers {
		serveMux.HandleFunc(h.Path, h.HandlerFunc)
	}
	s := &http.Server{Addr: ":" + port, Handler: serveMux}
	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Printf("HTTP server closing with message: %v", err)
		}
	}()
	log.Printf("[Start] HTTP server on port %s started\n", port)

	return s
}

func ShutdownHTTPServer(appName string, server *http.Server) {
	log.Printf("[Shutdown] %s HTTP server is shutting down\n", appName)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Unable to stop HTTP server: %v", err)
	}
}
