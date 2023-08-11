package main

import (
	"fmt"
	"net"
)

func main() {
	c, _ := net.Dial("tcp4", "127.0.0.1:8005")
	aa := make([]byte, 300)
	c.Read(aa)
	fmt.Println(aa)
}
