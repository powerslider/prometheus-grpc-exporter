package tcp

import (
	"io"
	"log"
	"net"

	"github.com/powerslider/prometheus-grpc-exporter/pkg/sd"
)

type Forwarder struct {
	Name            string
	LocalAddr       string
	RemoteAddresses []string
	HealtcheckAddr  string
}

func NewTCPForwarder(name string, localAddr string, remoteAddresses []string, healthcheckAddr string) *Forwarder {
	return &Forwarder{
		Name:            name,
		LocalAddr:       localAddr,
		RemoteAddresses: remoteAddresses,
		HealtcheckAddr:  healthcheckAddr,
	}
}

func (f *Forwarder) forward(serverConn net.Conn, remoteAddr string, consulService *sd.Consul) {
	client := NewTCPClient(remoteAddr)
	client.ConsulConnect(consulService, func(clientConn net.Conn) {
		log.Printf("Forwarding from %v to %v\n", serverConn.LocalAddr(), clientConn.RemoteAddr())
		go func() {
			defer clientConn.Close()
			defer serverConn.Close()
			io.Copy(clientConn, serverConn)
		}()
		go func() {
			defer clientConn.Close()
			defer serverConn.Close()
			io.Copy(serverConn, clientConn)
		}()
	})
}

func (f *Forwarder) Accept() {
	listener, err := NewTCPListener(f.LocalAddr)
	if err != nil {
		log.Fatal("Error creating TCP listener for TCP forwarder: ", err)
		return
	}
	server := NewTCPServer(f.Name, f.LocalAddr, listener)
	consulService, err := sd.NewConsulRegistration(f.Name, f.Name, f.HealtcheckAddr)
	if err != nil {
		log.Fatalf("Error registering %s with Consul server: %v", f.Name, err)
	}

	server.Accept(func(conn net.Conn) {
		for _, address := range f.RemoteAddresses {
			go f.forward(conn, address, consulService)
		}
	})
}
