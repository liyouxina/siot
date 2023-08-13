package main

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/liyouxina/siot/entity"
	log "github.com/sirupsen/logrus"
	"net"
	"strconv"
	"sync"
	"time"
)

var deviceIdAgentPool map[string]*Agent
var systemIdAgentPool map[string]*Agent

type Agent struct {
	Coon     net.Conn
	DeviceId string
	SystemId string
	Status   string
	mutex    sync.Mutex
}

const (
	HEX_GET_DEVICE_ID = "5aa506ffffffff00504c"
)

const (
	COMMAND_GET_DEVICE_ID = byte(50)
	COMMAND_GET_STATUS    = byte(51)
	COMMAND_TURN_ON_OFF   = byte(52)
)

const (
	COMMAND_GET_DEVICE_ID_STRING = "获取设备号"
	COMMAND_GET_STATUS_STRING    = "获取设备信息"
	COMMAND_TURN_ON_OFF_STRING   = "开关灯"
)

func byteServe() {
	listen, err := net.Listen("tcp4", "0.0.0.0:8001")
	if err != nil {
		fmt.Println("listen failed, err:", err)
		panic(err)
	}
	for {
		conn, err := listen.Accept() // 建立连接
		if err != nil {
			log.Warnf("接收连接出错 %s", err.Error())
			continue
		}
		agent := Agent{
			Coon:     conn,
			mutex:    sync.Mutex{},
			SystemId: strconv.FormatInt(time.Now().UnixMilli(), 10),
		}
		systemIdAgentPool[agent.SystemId] = &agent
		resp, err := agent.GetDeviceId()
		if err != nil {
			log.Warnf("接收连接 获取设备号出错 %s", err.Error())
			agent.Status = "接收连接 获取设备号出错"
		}
		agent.DeviceId = resp.DeviceId
		deviceIdAgentPool[resp.DeviceId] = &agent
		deviceDO := entity.GetByDeviceId(agent.DeviceId)
		if deviceDO == nil {
			deviceDO = &entity.Device{
				Name:     "lamp " + agent.SystemId,
				DeviceId: agent.DeviceId,
				Status:   agent.Status,
			}
			tx := entity.CreateDevice(deviceDO)
			if tx.Error != nil {
				log.Warnf("接受连接 设备写入数据库出错 %s", tx.Error.Error())
			}
			if tx.RowsAffected == 0 {
				log.Warnf("接受连接 设备写入数据库失败")
			}
		}
	}
}

func (agent *Agent) GetDeviceId() (*ByteResp, error) {
	agent.mutex.Lock()
	conn := agent.Coon
	byteArray, err := hex.DecodeString(HEX_GET_DEVICE_ID)
	if err != nil {
		log.Warnf("向设备请求设备号出错 十六进制转换出错 %s", err.Error())
		return nil, err
	}
	_, err = conn.Write(byteArray)
	if err != nil {
		log.Warnf("向设备请求设备号出错 发送数据出错 %s", err.Error())
		return nil, err
	}
	reader := bufio.NewReader(conn)
	var buf [4096]byte
	n, err := reader.Read(buf[:]) // 读取数据
	if err != nil {
		log.Warnf("向设备请求设备号出错 接收返回消息出错 %s", err.Error())
		return nil, err
	}

	log.WithField("接收数据", buf[:n]).Info("向设备请求设备号接收到数据")
	resp, err := getRespMsg(buf[:n])
	if err != nil {
		log.Warnf("向设备请求设备号 解析返回数据格式出错 %s", err.Error())
		return nil, err
	}
	agent.mutex.Unlock()
	log.WithField("接受数据", resp).Info("向设备请求设备号 返回成功")
	return resp, nil
}

func (agent *Agent) GetDeviceInfo() (*ByteResp, error) {
	return nil, nil
}

func (agent *Agent) sendMsg(content string) (*string, error) {
	agent.mutex.Lock()
	conn := agent.Coon
	byteArray, err := hex.DecodeString(content)
	if err != nil {
		log.Warnf("向设备发送消息出错 十六进制转换出错 %s", err.Error())
		return nil, err
	}
	_, err = conn.Write(byteArray)
	if err != nil {
		log.Warnf("向设备发送消息出错 发送数据出错 %s", err.Error())
		return nil, err
	}
	reader := bufio.NewReader(conn)
	var buf [4096]byte
	n, err := reader.Read(buf[:]) // 读取数据
	if err != nil {
		log.Warnf("向设备发送消息出错 接收返回消息出错 %s", err.Error())
		return nil, err
	}

	log.WithField("接收数据", buf[:n]).Info("向设备请求设备号接收到数据")
	resp := string(buf[:n])
	return &resp, nil
}

type ByteResp struct {
	Type     string
	DeviceId string
	Length   int
}

func getRespMsg(content []byte) (*ByteResp, error) {
	if len(content) < 11 {
		return nil, errors.New("解析包 返回包的长度不够")
	}
	if content[0] != byte(0x5a) || content[1] != byte(0xa5) {
		return nil, errors.New("解析包 包头部不正确")
	}
	length := int(content[2])
	if length > len(content)-4 {
		return nil, errors.New("解析包 返回包长度大于实际包长度")
	}
	result := &ByteResp{
		Length:   length,
		Type:     getCommand(content[8]),
		DeviceId: hex.EncodeToString(content[3:7]),
	}
	if result.Type == "" {
		return nil, errors.New("解析包 返回未知的操作指令")
	}
	return result, nil
}

func getCommand(command byte) string {
	if command == COMMAND_GET_DEVICE_ID {
		return COMMAND_GET_STATUS_STRING
	} else if command == COMMAND_GET_STATUS {
		return COMMAND_GET_STATUS_STRING
	} else if command == COMMAND_TURN_ON_OFF {
		return COMMAND_TURN_ON_OFF_STRING
	} else {
		return ""
	}
}

func genVerifyCode(verifyContent []byte) byte {
	a := byte(0)
	for _, c := range verifyContent {
		a = a + c
	}
	return a
}
