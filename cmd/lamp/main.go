package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/liyouxina/siot/entity"
	"net"
	"time"
)

func main() {
	serve()
	byteServe()
	go monitor()
}

func monitor() {
	for {
		time.Sleep(time.Minute)

	}
}

func scanDevices() {
	for {
		cursor := int64(0)
		devices := entity.ListDeviceByCursor(cursor, 50)
		if devices == nil || len(devices) == 0 {
			break
		}
		cursor = devices[len(devices)-1].Id
		for _, device := range devices {
			keepDevice(device)
		}
	}

}

func keepDevice(device *entity.Device) {
	conn := agentPool[device.DeviceId]
	if conn == nil {
		// mysql断开链接状态
		entity.ChangeDeviceState(device.Id, "disconnected")
	}

}

func serve() {
	server := gin.Default()
	server.GET("/command", command)
	server.GET("", command)
	_ = server.Run("0.0.0.0:8002")
}

const (
	GET_INFO   = "getInfo"
	OPEN_LAMP  = "open"
	CLOSE_LAMP = "close"
)

type Resp struct {
	Msg string `json:"msg"`
}

func command(context *gin.Context) {
	cmd := context.Query("command")
	deviceId := context.Query("deviceId")
	agent := agentPool[deviceId]
	if agent == nil {
		context.JSON(200, Resp{
			Msg: "没找到这个设备",
		})
		return
	}
	if GET_INFO == cmd {

	} else if OPEN_LAMP == cmd {

	} else if CLOSE_LAMP == cmd {

	} else {
		context.JSON(200, Resp{
			Msg: "没有这个操作",
		})
	}
}

func byteServe() {
	listen, err := net.Listen("tcp", "0.0.0.0:8001")
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
		getDeviceIdCmdBytes, _ := hex.DecodeString(HEX_GET_DEVICE_ID)
		_, _ = conn.Write(getDeviceIdCmdBytes)
		var buf [4096]byte
		reader := bufio.NewReader(conn)
		n, _ := reader.Read(buf[:])
		recvStr := string(buf[:n])
		agentPool[recvStr] = &Agent{
			Coon: conn,
		}
		go process(conn) // 启动一个goroutine处理连接
	}
}

const (
	HEX_GET_DEVICE_ID = "50"
	HEX_GET_STATUS    = "51"
	HEX_TURN_ON_OFF   = "52"
)

var agentPool map[string]*Agent

type Agent struct {
	Coon   net.Conn
	Status string
}

func process(conn net.Conn) {
	defer conn.Close()
	for {
		reader := bufio.NewReader(conn)
		var buf [4096]byte
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
