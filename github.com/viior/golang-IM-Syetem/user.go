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
	//查询在线用户
	case msg == "who":
		user.server.mapLock.Lock()
		for _, MapUser := range user.server.OnlineMap {
			onLineMsg := "[" + MapUser.Addr + "]" + MapUser.Name + " online\n"
			user.sendMsg(onLineMsg)
		}
		user.server.mapLock.Unlock()
		//重命名的分支
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
	case len(msg) > 4 && strings.Contains(msg, "->"):
		//消息格式：消息体->用户名
		//确认用户名
		remoteName := strings.Split(msg, "->")[1]
		if remoteName == "" {
			user.sendMsg("The message format is incorrect, please use the format of message body->username.\n")
			return
		}
		//查找用户
		if remoteUser, ok := user.server.OnlineMap[remoteName]; !ok {
			user.sendMsg("The target user does not exist\n")
			return
		} else {
			//向用户发送消息
			contents := strings.Split(msg, "->")[0]
			if contents == "" {
				user.sendMsg("No message content, please resend\n")
				return
			}
			remoteUser.sendMsg(user.Name + "send to you : " + contents + "\n")
		}
		//默认
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
