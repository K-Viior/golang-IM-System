package main

import (
	"fmt"
	"log"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int
	//在线用户列表
	OnlineMap map[string]*User
	mapLock   sync.Mutex
	//消息广播的Channel
	Message chan string
}

// 创建一个Server对象
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

// 处理业务的接口
func (server *Server) Handle(conn net.Conn) {
	//当前链接的业务
	//fmt.Println("链接建立成功")
	user := NewUser(conn)

	//将上线的用户加入在线列表中
	server.mapLock.Lock()
	server.OnlineMap[user.Name] = user
	server.mapLock.Unlock()
	//广播用户
	server.BroadCast(user, " is Online \n")

	//使当前Handler阻塞
	select {}
}

// 广播消息的方法
func (server *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + msg
	server.Message <- sendMsg
}

// 监听Message，一旦有消息，就广播给每个用户
func (server *Server) ListenMessage() {
	go func() {
		for {
			msg := <-server.Message
			//将message发送给每个在线的User
			server.mapLock.Lock()
			for _, cil := range server.OnlineMap {
				cil.C <- msg
			}
			server.mapLock.Unlock()
		}
	}()
}

// 启动服务器的接口
func (server *Server) Start() {
	//监听端口
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.Ip, server.Port))
	if err != nil {
		log.Printf("err happend in listen, err:", err)
		return
	}
	//关闭监听端口
	defer listen.Close()
	//开启监听服务端channel
	server.ListenMessage()
	for {
		//监听成功
		conn, err := listen.Accept()
		if err != nil {
			log.Printf("err happend in accept, err:", err)
			continue
		}
		//处理业务
		go server.Handle(conn)
	}

}
