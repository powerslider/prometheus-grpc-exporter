package tcp

import (
	"fmt"
	"log"
	"net"
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
	for {
		clientProcessingFunc(conn)
	}
}
