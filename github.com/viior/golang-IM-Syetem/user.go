package main

import (
	"net"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

// 用户上线的接口
func (user *User) OnLine() {
	//将上线的用户加入在线列表中
	user.server.mapLock.Lock()
	user.server.OnlineMap[user.Name] = user
	user.server.mapLock.Unlock()
	//广播用户
	user.server.BroadCast(user, " is Online \n")
}

// 用户下线的接口
func (user *User) OffLine() {
	//将上线的用户从在线列表中删除
	user.server.mapLock.Lock()
	delete(user.server.OnlineMap, user.Name)
	user.server.mapLock.Unlock()
	//广播用户
	user.server.BroadCast(user, " is Offline \n")
}

// 用户发送消息的接口
func (user *User) DoMessage(msg string) {
	user.server.BroadCast(user, msg)
}
func (user *User) ListenMsg() {
	go func() {
		for {
			msg := <-user.C
			user.conn.Write([]byte(msg))
		}
	}()
}

func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}
	//开启监听channel
	user.ListenMsg()
	return &user
}
