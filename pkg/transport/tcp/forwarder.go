package tcp

import (
	"io"
	"log"
	"net"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/connect"
)

type Forwarder struct {
	Name       string
	LocalAddr  string
	RemoteAddr string
}

func NewTCPForwarder(name string, localAddr string, remoteAddr string) *Forwarder {
	return &Forwarder{
		Name:       name,
		LocalAddr:  localAddr,
		RemoteAddr: remoteAddr,
	}
}

func (f *Forwarder) forward(serverConn net.Conn, consulService *connect.Service, consulClient *api.Client) {
	client := NewTCPClient(f.RemoteAddr)
	client.ConsulConnect(consulService, consulClient, func(clientConn net.Conn) {
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
	// Create a Consul API client
	consulClient, _ := api.NewClient(api.DefaultConfig())

	consulService, _ := connect.NewService(f.Name, consulClient)
	defer consulService.Close()

	server.Accept(func(conn net.Conn) {
		go f.forward(conn, consulService, consulClient)
	})
}
