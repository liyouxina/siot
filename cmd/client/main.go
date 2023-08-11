package main

import (
	"fmt"
	"net"
	"time"
)

func sender(conn net.Conn) {
	words := "hello world!"
	conn.Write([]byte(words))
	buf := make([]byte, 4096)
	conn.Read(buf)
	fmt.Println(buf)
	fmt.Println("send over")

}

var agentPool map[string]*Agent

type Agent struct {
	Coon   net.Conn
	Status string
}

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8001")
	if err != nil {
		fmt.Println(err.Error())
	}
	for {
		time.Sleep(time.Millisecond * 300)
		conn.Write([]byte{2})
	}
}
