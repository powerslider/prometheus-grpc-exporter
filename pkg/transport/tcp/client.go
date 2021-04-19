package tcp

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/hashicorp/consul/api"

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

func (c *Client) ConsulConnect(consulService *connect.Service, consulClient *api.Client, clientProcessingFunc ClientProcessor) {
	conn, err := consulService.Dial(context.Background(), &connect.ConsulResolver{
		Client: consulClient,
		Name:   c.ServerAddr,
	})
	log.Printf("Dialed %s", c.ServerAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	clientProcessingFunc(conn)
}
