package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	commons "rpcChat/commons"
)

type Client commons.Client
type Message commons.Message

func main() {
	rpcClient, err := rpc.Dial("tcp", "localhost:42586")
	if err != nil {
		log.Fatal(err)
	}

	host := "localhost"
	port, _ := GetFreePort()
	addr := host + ":" + fmt.Sprint(port)
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	print("Welcome\033[0m to \033[36mRPC-Chat\033[0m!\n\nType \033[31m exit \033[0m to leave the chat\n\n")
	print("Choose a handle: ")
	reader := bufio.NewScanner(os.Stdin)
	reader.Scan()
	name := reader.Text()
	cli := Client{Addr: addr, Name: name}
	err = rpcClient.Call("Client.ADD", cli, new(bool))

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		inbound, err := net.ListenTCP("tcp", tcpAddr)
		if err != nil {
			log.Fatal(err)
		}
		listener := new(Client)
		rpc.Register(listener)
		rpc.Accept(inbound)
	}()

	defer func() {
		println("Bye!")
		rpcClient.Call("Client.REMOVE", cli, new(bool))
	}()
	
	in := bufio.NewReader(os.Stdin)
	for {
		msg, _, err := in.ReadLine()
		if string(msg) == "exit" {
			return
		}
		if err != nil {
			log.Fatal(err)
		}
		envelop := commons.Message{Sender: commons.Client(cli), Content: string(msg)}
		err = rpcClient.Call("Client.HandleMessage", envelop, new(bool))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (c *Client) Listen(msg string, reply *bool) error {
	*reply = true
	print(msg + "\n\033[34m>\033[0m ")
	return nil
}

func GetFreePort() (port int, err error) {
	var a *net.TCPAddr
	if a, err = net.ResolveTCPAddr("tcp", "localhost:0"); err == nil {
		var l *net.TCPListener
		if l, err = net.ListenTCP("tcp", a); err == nil {
			defer l.Close()
			return l.Addr().(*net.TCPAddr).Port, nil
		}
	}
	return
}
