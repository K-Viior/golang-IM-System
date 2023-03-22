package main

import (
	"net"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn
}

func (user *User) ListenMsg() {
	go func() {
		for {
			msg := <-user.C
			user.conn.Write([]byte(msg))
		}
	}()
}

func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()
	user := User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,
	}
	//开启监听channel
	user.ListenMsg()
	return &user
}
