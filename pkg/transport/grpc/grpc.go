package grpc

import (
	"context"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

type Server struct {
	Name           string
	Port           string
	Listener       net.Listener
	serverInstance *grpc.Server
}

func NewGRPCServer(name string, port string, listener net.Listener) *Server {
	return &Server{
		Name:           name,
		Port:           port,
		Listener:       listener,
		serverInstance: grpc.NewServer(),
	}
}

func (s *Server) Start(serviceServerRegistrarFunc func(*grpc.Server)) {
	serviceServerRegistrarFunc(s.serverInstance)
	go func() {
		if err := s.serverInstance.Serve(s.Listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	log.Printf("[Start] %s gRPC server on port %s started\n", s.Name, s.Port)
}

func (s *Server) Shutdown() {
	log.Printf("[Shutdown] %s gRPC server is shutting down\n", s.Name)

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if s.serverInstance != nil {
		s.serverInstance.GracefulStop()
	}
}
