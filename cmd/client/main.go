package main

import (
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
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
	aa := []byte{0xff, 0xff, 0xff, 0xff, 0x00, 0x50}
	var a byte
	a = byte(0)
	for _, v := range aa {
		a = a + v
	}
	aaaa, _ := hex.DecodeString(strconv.Itoa(int(a)))
	fmt.Println(aaaa)
}
