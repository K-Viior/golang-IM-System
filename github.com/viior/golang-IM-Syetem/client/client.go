package main

import (
	"fmt"
	"log"
	"net"
)

type Client struct {
	ServerIp string
	Port     int
	Name     string
	conn     net.Conn
}

func NewClient(serverIp string, port int) *Client {
	//创建客户端
	client := &Client{
		ServerIp: serverIp,
		Port:     port,
	}
	//进行链接
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, port))
	if err != nil {
		log.Printf("Error in dial : ", err)
		return nil
	}

	client.conn = conn
	//返回客户端
	return client
}

func main() {
	client := NewClient("127.0.0.1", 8888)
	if client == nil {
		log.Printf(">>>>>>>Server connect failed>>>>>")
		return
	}
	fmt.Println(">>>>>>>>Server connect success>>>>>>")
	select {}
}
