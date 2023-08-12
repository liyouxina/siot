package main

import "net"

func main() {
	listen, _ := net.Listen("tcp4", "0.0.0.0:8005")
	conn, _ := listen.Accept()
	conn.Write([]byte{1, 2, 3})
	conn.Close()
}
