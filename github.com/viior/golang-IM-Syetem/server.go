package main

import (
	"fmt"
	"log"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

// 创建一个Server对象
func NewServer(ip string, port int) *Server {
	server := Server{
		Ip:   ip,
		Port: port,
	}
	return &server
}
func (server *Server) Handle(conn net.Conn) {
	//当前链接的业务
	fmt.Println("链接建立成功")
}

// 启动服务器的接口
func (server *Server) Start() {
	//监听端口
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.Ip, server.Port))
	if err != nil {
		log.Printf("err happend in listen, err:", err)
		return
	}
	//关闭监听
	defer listen.Close()

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
