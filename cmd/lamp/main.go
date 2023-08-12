package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/liyouxina/siot/entity"
	"net"
	"strconv"
	"strings"
	"time"
)

func main() {
	agentPool = map[string]*Agent{}
	systemIdAgentPool = map[string]*Agent{}
	go serve()
	byteServe()
	go monitor()
	go releaseDisconnectedAgent()
}

func releaseDisconnectedAgent() {
	for {
		time.Sleep(time.Second)
		for k, agent := range agentPool {
			if isConnected(agent.Coon) {
				_ = agent.Coon.Close()
				agentPool[k] = nil
				systemIdAgentPool[k] = nil
			}
		}
		for k, agent := range systemIdAgentPool {
			if isConnected(agent.Coon) {
				_ = agent.Coon.Close()
				agentPool[k] = nil
				systemIdAgentPool[k] = nil
			}
		}
	}

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
	GET_ALL_AGENTS        = "getAllAgents"
	GET_ALL_BY_DEVICE_ID  = "getAllAgentsByDeviceId"
	GET_INFO              = "getInfo"
	OPEN_LAMP             = "open"
	CLOSE_LAMP            = "close"
	SEND_MSG_BY_SYSTEM_ID = "sendMsgBySystemId"
	SEND_MSG_BY_ID        = "sendMsgById"
)

type Resp struct {
	Msg string `json:"msg"`
}

func command(context *gin.Context) {
	cmd := context.Query("command")
	// deviceId := context.Query("deviceId")
	systemId := context.Query("systemId")
	hexContent := context.Query("hex")

	if GET_INFO == cmd {

	} else if OPEN_LAMP == cmd {

	} else if CLOSE_LAMP == cmd {

	} else if GET_ALL_BY_DEVICE_ID == cmd {
		context.JSON(200, fmt.Sprintln("%v", agentPool))
	} else if SEND_MSG_BY_SYSTEM_ID == cmd {
		agent := systemIdAgentPool[systemId]
		if agent == nil {
			context.JSON(200, Resp{
				Msg: "没有这个设备",
			})
			return
		}
		hexContent = strings.Join(strings.Split(hexContent, " "), "")
		content, err := hex.DecodeString(hexContent)
		if err != nil {
			context.JSON(200, Resp{
				Msg: err.Error(),
			})
			return
		}
		_, err = agent.Coon.Write(content)
		if err != nil {
			context.JSON(200, Resp{
				Msg: err.Error(),
			})
			return
		}
		reader := bufio.NewReader(agent.Coon)
		var buf [4096]byte
		n, err := reader.Read(buf[:]) // 读取数据
		if err != nil {
			context.JSON(200, Resp{
				Msg: "接收数据有问题" + err.Error(),
			})
			return
		}
		recvStr := string(buf[:n])
		context.JSON(200, Resp{
			Msg: fmt.Sprintf("发送成功 返回值 %s", recvStr),
		})

	} else if GET_ALL_AGENTS == cmd {
		context.JSON(200, fmt.Sprintln("%v", systemIdAgentPool))
	} else {
		context.JSON(200, Resp{
			Msg: "没有这个操作",
		})
	}
}

func byteServe() {
	listen, err := net.Listen("tcp4", "0.0.0.0:8001")
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
		systemIdAgentPool[strconv.Itoa(time.Now().Minute())] = &Agent{
			Coon: conn,
		}
		go registerDeviceId(conn)
		go process(conn)
	}
}

func registerDeviceId(conn net.Conn) {
	getDeviceIdCmdBytes, _ := hex.DecodeString(HEX_GET_DEVICE_ID)
	_, _ = conn.Write(getDeviceIdCmdBytes)
	var buf [4096]byte
	reader := bufio.NewReader(conn)
	n, _ := reader.Read(buf[:])
	recvStr := string(buf[:n])
	agentPool[recvStr] = &Agent{
		Coon: conn,
	}
}

const (
	HEX_GET_DEVICE_ID = "50"
	HEX_GET_STATUS    = "51"
	HEX_TURN_ON_OFF   = "52"
)

var agentPool map[string]*Agent
var systemIdAgentPool map[string]*Agent

type Agent struct {
	Coon   net.Conn
	Status string
}
