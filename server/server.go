package main

import (
	"log"
	"net"
	"net/rpc"
	commons "rpcChat/commons"
)

var (
	messages = make(chan string, 50)
	clients  = make(map[User]bool)
)

type User commons.User
type Message commons.Message

func main() {
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:42586")
	if err != nil {
		log.Fatal(err)
	}

	inbound, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	listener := new(User)
	rpc.Register(listener)
	go rpc.Accept(inbound)

	for msg := range messages {
		for cli := range clients {
			inb, err := rpc.Dial("tcp", cli.Addr)
			if err != nil {
				continue
			}
			err = inb.Call("User.Listen", msg, new(bool))
			if err != nil {
				continue
			}
		}
	}
}

func (c *User) ADD(C User, reply *bool) error {
	clients[C] = true
	println("registered user " + C.Name + " with adress " + "[" + C.Addr + "]")
	messages <- "[" + C.Name + "]" + " has joined the chat!" + "\n"
	return nil
}

func (c *User) ProcessMessage(msg Message, reply *bool) error {
	messages <-  msg.Sender.Name + " : " + msg.Content + "\n"
	return nil
}
