package main

import (
	"log"
	"net"
	"net/rpc"
	commons "rpcChat/commons"
)

type Client commons.Client
type Message commons.Message

var (
	messages = make(chan string, 500)
	clients  = make(map[Client]bool)

	ADDS = []string{
		"Adds By Creator: Bayern Rules",
		"Adds By Creator: Real 's Lame",
		"Adds By Creator: LEVA 's GOAT",
		"Adds By Creator: Barca My A*$",
	}
	ADDAGRESS = 15
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
		}
	}
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
