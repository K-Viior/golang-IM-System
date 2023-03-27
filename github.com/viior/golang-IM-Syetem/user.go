package main

import (
	"net"
	"strings"
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
	user.server.BroadCast(user, " Sign in ") //is Online
}

// 用户下线的接口
func (user *User) OffLine() {
	//将上线的用户从在线列表中删除
	user.server.mapLock.Lock()
	delete(user.server.OnlineMap, user.Name)
	user.server.mapLock.Unlock()
	//广播用户
	user.server.BroadCast(user, " Sign out ")
}

// 用户发送消息的接口
func (user *User) DoMessage(msg string) {
	switch {
	case msg == "who":
		//查询在线用户
		user.server.mapLock.Lock()
		for _, MapUser := range user.server.OnlineMap {
			onLineMsg := "[" + MapUser.Addr + "]" + MapUser.Name + " online\n"
			user.sendMsg(onLineMsg)
		}
		user.server.mapLock.Unlock()
	case len(msg) > 7 && msg[:7] == "rename|":
		//用户名重命名
		newName := strings.Split(msg, "|")[1]
		//判断新名称是否存在
		_, ok := user.server.OnlineMap[newName]
		if ok {
			user.sendMsg("This name has been used\n")
		} else {
			user.server.mapLock.Lock()
			delete(user.server.OnlineMap, user.Name)
			user.server.OnlineMap[newName] = user
			user.server.mapLock.Unlock()
			user.Name = newName
			user.sendMsg("Your name " + newName + " update succeeded\n")
		}
	default:
		user.server.BroadCast(user, msg)
	}
}
func (user *User) ListenMsg() {
	go func() {
		for {
			msg := <-user.C
			user.sendMsg(msg)
		}
	}()
}

func (user *User) sendMsg(msg string) {
	user.conn.Write([]byte(msg))
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
