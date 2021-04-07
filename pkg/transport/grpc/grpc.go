package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

func StartGRPCServer(port string, serviceServerRegistrarFunc func(*grpc.Server)) *grpc.Server {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	serviceServerRegistrarFunc(s)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	log.Printf("[Start] gRPC server on port %s started\n", port)

	return s
}

func ShutdownGRPCServer(appName string, server *grpc.Server) {
	log.Printf("[Shutdown] %s gRPC server is shutting down\n", appName)

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if server != nil {
		server.GracefulStop()
	}
}
