package tcp

import (
	"fmt"
	"log"
	"net"
)

type ServerProcessor func(conn net.Conn)

type server struct {
	port string
}

func NewTCPServer(port string) *server {
	return &server{
		port: port,
	}
}

func (s *server) Accept(serverProcessingFunc ServerProcessor) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", s.port))
	if err != nil {
		log.Fatalf("Failed to setup listener: %v", err)
	}
	log.Printf("[Start] TCP server on port %s started\n", s.port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("ERROR: failed to accept listener: %v", err)
		}
		log.Printf("Accepted connection from %v\n", conn.RemoteAddr().String())
		serverProcessingFunc(conn)
	}

}
