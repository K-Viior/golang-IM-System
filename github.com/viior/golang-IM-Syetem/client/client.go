package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
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

// 更新用户名的方法
func (client *Client) UpdateName() bool {
	fmt.Println("Please enter a new name")
	fmt.Scanln(&client.Name)
	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.write err : ", err)
		return false
	}
	return true
}

// 公聊方法
func (client *Client) PublicChat() {
	var chatMsg string
	fmt.Println(">>>>Please enter the message ,exit will be down<<<<")
	fmt.Scanln(&chatMsg)
	for chatMsg != "exit" {
		//发送给服务器
		if len(chatMsg) > 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				log.Printf("Conn write error : ", err)
				break
			}
		}
		chatMsg = ""
		fmt.Println(">>>>Please enter the message ,exit will be down<<<<")
		fmt.Scanln(&chatMsg)
	}
}

// 私聊方法
func (client *Client) PrivateChat() {
	var remoteName string
	var chatMsg string

	client.SelectName()
	fmt.Println(">>>>Please enter the user who you want chat,exit will be down<<<<")
	fmt.Scanln(&remoteName)
	for remoteName != "exit" {
		fmt.Println(">>>>Please enter the message ,exit will be down<<<<")
		fmt.Scanln(&chatMsg)
		for chatMsg != "exit" {
			if len(chatMsg) > 0 {
				sendMsg := chatMsg + "->" + remoteName + "\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					log.Printf("Conn write error : ", err)
					break
				}
			}
			chatMsg = ""
			fmt.Println(">>>>Please enter the message ,exit will be down<<<<")
			fmt.Scanln(&chatMsg)
		}
		client.SelectName()
		fmt.Println(">>>>Please enter the user who you want chat,exit will be down<<<<")
		fmt.Scanln(&remoteName)
	}
}
func (client *Client) SelectName() {
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		log.Printf("Conn write error : ", err)
		return
	}
}
func (client *Client) run() {
	for client.flag != 0 {
		for client.menu() != true {
		}
		switch client.flag {
		case 1:
			client.PublicChat()
		case 2:
			client.PrivateChat()
		case 3:
			client.UpdateName()
		}
	}
}

var ServerIp string
var Port int

func init() {
	flag.StringVar(&ServerIp, "ip", "127.0.0.1", "Set the server IP address (default is 127.0.0.1)")
	flag.IntVar(&Port, "port", 8888, "Set the server port (default is 8888).")
}

// 处理server回应的消息，直接显示标准输出
func (client *Client) DealResponse() {
	io.Copy(os.Stdout, client.conn)
}
func main() {
	//命令行解析
	flag.Parse()
	client := NewClient(ServerIp, Port)
	if client == nil {
		log.Printf(">>>>>>>Server connect failed>>>>>")
		return
	}
	//goroutine 开启处理回执消息
	go client.DealResponse()
	fmt.Println(">>>>>>>>Server connect success>>>>>>")
	client.run()
}
