package main

import (
	"flag"
	"fmt"
	"log"
	"net"
)

type Client struct {
	ServerIp string
	Port     int
	Name     string
	conn     net.Conn
	flag     int
}

func NewClient(serverIp string, port int) *Client {
	//创建客户端
	client := &Client{
		ServerIp: serverIp,
		Port:     port,
		flag:     999,
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
func (client *Client) menu() bool {
	var flag int
	fmt.Println("1.public chat mode")
	fmt.Println("2.private chat mode")
	fmt.Println("3.update username")
	fmt.Println("0.exit")
	fmt.Scanln(&flag)
	if flag > 3 || flag < 0 {
		fmt.Println("Please enter a number within the valid range")
		return false
	}
	client.flag = flag
	return true
}

func (client *Client) run() {
	for client.flag != 0 {
		for client.menu() != true {
		}
		switch client.flag {
		case 1:
			fmt.Println("public chat mode")
		case 2:
			fmt.Println("private chat mode")
		case 3:
			fmt.Println("update username")
		}
	}
}

var ServerIp string
var Port int

func init() {
	flag.StringVar(&ServerIp, "ip", "127.0.0.1", "Set the server IP address (default is 127.0.0.1)")
	flag.IntVar(&Port, "port", 8888, "Set the server port (default is 8888).")
}

func main() {
	//命令行解析
	flag.Parse()
	client := NewClient(ServerIp, Port)
	if client == nil {
		log.Printf(">>>>>>>Server connect failed>>>>>")
		return
	}
	fmt.Println(">>>>>>>>Server connect success>>>>>>")
	client.run()
}
