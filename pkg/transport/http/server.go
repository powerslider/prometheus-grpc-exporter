package http

import (
	"context"
	"log"
	"net/http"
	"time"
)

type ServerWrapper struct {
	Name           string
	Port           string
	serverInstance *http.Server
}

type Handler struct {
	HandlerFunc func(http.ResponseWriter, *http.Request)
	Path        string
}

func NewHTTPServer(name string, port string, handlers ...Handler) *ServerWrapper {
	serveMux := http.NewServeMux()
	for _, h := range handlers {
		serveMux.HandleFunc(h.Path, h.HandlerFunc)
	}
	return &ServerWrapper{
		Name: name,
		Port: port,
		serverInstance: &http.Server{
			Addr:    ":" + port,
			Handler: serveMux,
		}}
}
func (s *ServerWrapper) Start() {
	go func() {
		if err := s.serverInstance.ListenAndServe(); err != nil {
			log.Printf("HTTP server closing with message: %v", err)
		}
	}()
	log.Printf("[Start] HTTP server on port %s started\n", s.Port)
}

func (s *ServerWrapper) Shutdown() {
	log.Printf("[Shutdown] %s HTTP server is shutting down\n", s.Name)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.serverInstance.Shutdown(ctx); err != nil {
		log.Fatalf("Unable to stop HTTP server: %v", err)
	}
}
