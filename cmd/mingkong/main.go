package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"net"
)

func process(conn net.Conn) {
	defer conn.Close() // 关闭连接
	for {
		reader := bufio.NewReader(conn)
		var buf [128]byte
		n, err := reader.Read(buf[:]) // 读取数据
		if err != nil {
			fmt.Println("read from client failed, err:", err)
			break
		}
		recvStr := string(buf[:n])
		fmt.Println("收到client端发来的数据：", recvStr)
		hexString := "a55aff001b4a53534831393035303030355d19bbde0000a32b55aa"
		byteArray, err := hex.DecodeString(hexString)
		if err != nil {
			fmt.Println("send failed, err:", err)
			break
		}
		conn.Write(byteArray) // 发送数据
	}
}

func main() {
	listen, err := net.Listen("tcp", "0.0.0.0:20000")
	if err != nil {
		fmt.Println("listen failed, err:", err)
		return
	}
	for {
		conn, err := listen.Accept() // 建立连接
		if err != nil {
			fmt.Println("accept failed, err:", err)
			continue
		}
		go process(conn) // 启动一个goroutine处理连接
	}
}
