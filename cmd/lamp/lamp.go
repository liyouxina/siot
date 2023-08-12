package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"net"
	"time"
)

func process(conn net.Conn) {
	defer conn.Close()
	for {
		byteArray, err := hex.DecodeString()
		if err != nil {
			fmt.Println("send failed, err:", err)
			break
		}
		_, _ = conn.Write(byteArray)
		reader := bufio.NewReader(conn)
		var buf [4096]byte
		n, err := reader.Read(buf[:]) // 读取数据
		if err != nil {
			fmt.Println("read from client failed, err:", err)
			break
		}
		recvStr := string(buf[:n])
		fmt.Println("收到client端发来的数据：", recvStr)
	}
}

func isConnected(conn net.Conn) bool {
	// 设置超时时间，以避免长时间等待
	conn.SetReadDeadline(time.Now())
	buffer := make([]byte, 1)

	_, err := conn.Read(buffer)
	if err != nil {
		if err, ok := err.(net.Error); ok && err.Timeout() {
			return true // 仍然连接着
		}
		return false // 连接已关闭或发生其他错误
	}

	return true // 仍然连接着
}
