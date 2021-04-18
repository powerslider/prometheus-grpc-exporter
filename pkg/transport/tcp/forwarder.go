package tcp

import (
	"io"
	"log"
	"net"
)

type Forwarder struct {
	LocalAddr  string
	RemoteAddr string
}

func NewTCPForwarder(localAddr string, remoteAddr string) *Forwarder {
	return &Forwarder{
		LocalAddr:  localAddr,
		RemoteAddr: remoteAddr,
	}
}

func (f *Forwarder) forward(conn net.Conn) {
	client, err := net.Dial("tcp", f.RemoteAddr)
	if err != nil {
		log.Printf("Dial failed: %v", err)
		defer conn.Close()
		return
	}
	log.Printf("Forwarding from %v to %v\n", conn.LocalAddr(), client.RemoteAddr())
	go func() {
		defer client.Close()
		defer conn.Close()
		io.Copy(client, conn)
	}()
	go func() {
		defer client.Close()
		defer conn.Close()
		io.Copy(conn, client)
	}()
}

func (f *Forwarder) Accept() {
	server := NewTCPServer(f.LocalAddr)
	server.Accept(func(conn net.Conn) {
		go f.forward(conn)
	})
}
