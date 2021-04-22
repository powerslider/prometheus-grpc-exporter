package tcp

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/powerslider/prometheus-grpc-exporter/pkg/sd"

	"github.com/hashicorp/consul/connect"
)

type ClientProcessor func(conn net.Conn)

type Client struct {
	ServerAddr string
}

func NewTCPClient(serverAddr string) *Client {
	return &Client{ServerAddr: serverAddr}
}

func (c *Client) Connect(clientProcessingFunc ClientProcessor) {
	conn, err := net.Dial("tcp", c.ServerAddr)
	log.Printf("Dialed %s", c.ServerAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	clientProcessingFunc(conn)
}

func (c *Client) ConsulConnect(consul *sd.Consul, clientProcessingFunc ClientProcessor) {
	conn, err := consul.Service.Dial(context.Background(), &connect.ConsulResolver{
		Client: consul.Client,
		Name:   c.ServerAddr,
	})
	log.Printf("Dialed %s", c.ServerAddr)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	clientProcessingFunc(conn)
}
