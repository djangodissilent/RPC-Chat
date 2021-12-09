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

type User commons.User
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

	fmt.Printf("Give me your name: ")
	reader := bufio.NewScanner(os.Stdin)
	reader.Scan()
	name := reader.Text()
	cli := User{Addr: addr, Name: name}
	err = rpcClient.Call("User.ADD", cli, new(bool))

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		inbound, err := net.ListenTCP("tcp", tcpAddr)
		if err != nil {
			log.Fatal(err)
		}
		listener := new(User)
		rpc.Register(listener)
		rpc.Accept(inbound)
	}()
	in := bufio.NewReader(os.Stdin)
	for {
		msg, _, err := in.ReadLine()
		if err != nil {
			log.Fatal(err)
		}
		envelop := commons.Message{Sender: commons.User(cli), Content: string(msg)}
		err = rpcClient.Call("User.ProcessMessage", envelop, new(bool))
		if err != nil {
			log.Fatal(err)
		}
	}
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

func (c *User) Listen(msg string, reply *bool) error {
	*reply = true
	print(msg + "\n> ")
	return nil
}