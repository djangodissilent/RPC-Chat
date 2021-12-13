package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	commons "rpcChat/commons"
	"time"
)

type Client commons.Client
type Message commons.Message

var (
	messages = make(chan string, 200)
	clients  = make(map[Client]bool)

	Interval = 3
)

func main() {
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:42586")
	if err != nil {
		log.Fatal(err)
	}

	inbound, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	listener := new(Client)
	rpc.Register(listener)
	go rpc.Accept(inbound)
	go spinner()

	ticker := time.NewTicker(time.Duration(Interval) * time.Second)

	for {
		select {

		case msg := <-messages:
			for cli := range clients {
				inb, err := rpc.Dial("tcp", cli.Addr)
				if err != nil {
					continue
				}
				err = inb.Call("Client.Listen", msg, new(bool))
				if err != nil {
					continue
				}
			}
		case <-ticker.C:
			fmt.Printf(" \033[32m Online: \033[33m %d \033[34m | \033[31m Server-Time: %s \033[33m", len(clients), time.Now().Format("2006-01-02 15:04:05"))
		}
	}
}

func (c *Client) REMOVE(C Client, reply *bool) error {
	*reply = true
	_, ok := clients[C]
	if ok {
		delete(clients, C)
		println("removed user " + C.Name + " with adress " + "[" + C.Addr + "]")
		messages <- "[" + C.Name + "]" + " has left the chat!" + "\n"
	}
	return nil
}

func (c *Client) ADD(C Client, reply *bool) error {
	clients[C] = true
	println("registered user " + C.Name + " with adress " + "[" + C.Addr + "]")
	messages <- "[" + C.Name + "]" + " has joined the chat!" + "\n"
	return nil
}

func (c *Client) HandleMessage(msg Message, reply *bool) error {
	messages <- "[" + msg.Sender.Name + "]: " + msg.Content + "\n"
	return nil
}

func spinner() {
	for {
		for _, r := range `-\|/` {
			fmt.Printf("\r%c", r)

			time.Sleep(time.Millisecond * 70)
		}
	}
}
