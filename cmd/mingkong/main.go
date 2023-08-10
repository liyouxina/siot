package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/liyouxina/siot/config"
	"net"
	"strconv"
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
	defer func() {
		err := recover()
		fmt.Println(err)
	}()
	agentPool = map[string]*Agent{}
	run()
}

func serve() {
	server := gin.Default()
	server.GET("/command", func(context *gin.Context) {

	})

	server.Run("0.0.0.0:8002")
}

const (
	GET_INFO   = "getInfo"
	OPEN_LAMP  = "open"
	CLOSE_LAMP = "close"
)

func command(context *gin.Context) {
	cmd := context.Query("command")
	if GET_INFO == cmd {

	}
}

var agentPool map[string]*Agent

type Agent struct {
	Coon   net.Conn
	Status string
}

func run() {
	port := strconv.Itoa(config.Config.MingKongConfig.Port)
	listen, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		fmt.Println("listen failed, err:", err)
		panic(err)
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

const (
	HEX_GET_DEVICE_ID = "50"
	HEX_GET_STATUS    = "51"
	HEX_TURN_ON_OFF   = "52"
)
