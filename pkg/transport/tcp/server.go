package tcp

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
)

type ServerProcessor func(conn net.Conn)

type Server struct {
	Name     string
	Port     string
	Listener net.Listener
}

func NewTCPListener(port string, tlsConfig ...*tls.Config) (net.Listener, error) {
	port = fmt.Sprintf(":%s", port)
	var listener net.Listener
	var err error
	if len(tlsConfig) > 0 {
		listener, err = tls.Listen("tcp", port, tlsConfig[0])
	} else {
		listener, err = net.Listen("tcp", port)
	}
	if err != nil {
		return nil, err
	}
	return listener, nil
}

func NewTCPServer(name string, port string, listener net.Listener) *Server {
	return &Server{
		Name:     name,
		Port:     port,
		Listener: listener,
	}
}

func (s *Server) Accept(serverProcessingFunc ServerProcessor) {
	log.Printf("[Start] %s TCP server on port %s started\n", s.Name, s.Port)
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			log.Fatalf("ERROR: failed to accept listener: %v", err)
		}
		log.Printf("Accepted connection from %v\n", conn.RemoteAddr().String())
		serverProcessingFunc(conn)
	}
}
