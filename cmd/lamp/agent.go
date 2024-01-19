package main

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/liyouxina/siot/entity"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net"
	"strconv"
	"strings"
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
	HEX_GET_DEVICE_ID   = "5aa506ffffffff00504c"
	HEX_GET_DEVICE_INFO = "5aa506%s0051"
	HEX_OPEN_CLOSE      = "5aa508%s0052%s%s"
)

const (
	COMMAND_GET_DEVICE_ID = byte(0x50)
	COMMAND_GET_STATUS    = byte(0x51)
	COMMAND_TURN_ON_OFF   = byte(0x52)
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
		log.Infof("接收连接 接收到连接")
		agent := Agent{
			Coon:     conn,
			mutex:    sync.Mutex{},
			SystemId: strconv.FormatInt(time.Now().UnixMilli(), 10),
		}
		systemIdAgentPool[agent.SystemId] = &agent
		resp, err := agent.GetDeviceId()
		if err != nil {
			log.Warnf("接收连接 获取设备号出错 %s %s", agent.SystemId, err.Error())
			agent.Status = "接收连接 获取设备号出错"
			continue
		}
		if resp == nil {
			log.Warnf("接收连接 获取设备号为空 %s", agent.SystemId)
			agent.Status = "接收连接 获取设备号为空"
			continue
		}
		agent.DeviceId = resp.DeviceId
		deviceIdAgentPool[resp.DeviceId] = &agent
		deviceDO, err := entity.GetByDeviceId(agent.DeviceId)
		if err != nil && strings.Index(err.Error(), gorm.ErrRecordNotFound.Error()) > 0 {
			log.Warnf("接受连接 读取数据库出错 %s %s %s", agent.SystemId, agent.DeviceId, err.Error())
			continue
		}
		if deviceDO == nil || deviceDO.Id == 0 {
			deviceDO = &entity.Device{
				Name:     "lamp " + agent.SystemId,
				DeviceId: agent.DeviceId,
				Status:   agent.Status,
				Type:     entity.LAMP_TYPE,
			}
			log.Infof("接收链接 创建新路灯设备 %s %s %s", agent.SystemId, agent.DeviceId, agent.Status)
			tx := entity.CreateDevice(deviceDO)
			if tx.Error != nil {
				log.Warnf("接受连接 设备写入数据库出错 %s %s %s", agent.SystemId, agent.DeviceId, tx.Error.Error())
			}
			if tx.RowsAffected == 0 {
				log.Warnf("接受连接 设备写入数据库失败 %s %s", agent.SystemId, agent.DeviceId)
			}
			log.Infof("接受连接 设备写入数据库成功 %s %s", agent.SystemId, agent.DeviceId)
		}
	}
}

func (agent *Agent) GetDeviceId() (*ByteResp, error) {
	agent.mutex.Lock()
	defer agent.mutex.Unlock()
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
	log.WithField("接受数据", resp).Info("向设备请求设备号 返回成功")
	return resp, nil
}

func (agent *Agent) GetDeviceInfo() (*ByteResp, error) {
	agent.mutex.Lock()
	defer agent.mutex.Unlock()
	if agent.DeviceId == "" {
		return nil, errors.New("获取设备信息 agent设备号为空")
	}
	conn := agent.Coon
	requestString := fmt.Sprintf(HEX_GET_DEVICE_INFO, agent.DeviceId)
	requestHexByte, err := hex.DecodeString(requestString)
	if err != nil {
		log.Warnf("获取设备信息 请求转换成16进制出错 %s", err.Error())
		return nil, err
	}
	verifyCode := genVerifyCode(requestHexByte[3:])
	requestHexByte = append(requestHexByte, verifyCode)
	log.WithField("请求内容", requestHexByte).Infof("获取设备信息 请求内容")

	_, err = conn.Write(requestHexByte)
	if err != nil {
		log.Warnf("获取设备信息 发送数据出错 %s", err.Error())
		return nil, err
	}
	reader := bufio.NewReader(conn)
	var buf [4096]byte
	n, err := reader.Read(buf[:]) // 读取数据
	if err != nil {
		log.Warnf("获取设备信息 接收返回消息出错 %s", err.Error())
		return nil, err
	}

	log.WithField("接收数据", buf[:n]).Info("获取设备信息 接收到数据")
	resp, err := getRespMsg(buf[:n])
	if err != nil {
		log.Warnf("获取设备信息 解析返回数据格式出错 %s", err.Error())
		return nil, err
	}
	log.WithField("接受数据", resp).Info("获取设备信息 返回成功")
	return resp, nil
}

func (agent *Agent) OpenClose(isOpen bool, roadId int) (*ByteResp, error) {
	agent.mutex.Lock()
	defer agent.mutex.Unlock()
	if agent.DeviceId == "" {
		return nil, errors.New("开关设备 agent设备号为空")
	}
	conn := agent.Coon
	commandString := ""
	if isOpen {
		commandString = fmt.Sprintf("0%d00", roadId)
	} else {
		commandString = fmt.Sprintf("000%d", roadId)
	}
	requestString := fmt.Sprintf(HEX_OPEN_CLOSE, agent.DeviceId, commandString)
	requestHexByte, err := hex.DecodeString(requestString)
	if err != nil {
		log.Warnf("开关设备 请求转换成16进制出错 %s", err.Error())
		return nil, err
	}
	verifyCode := genVerifyCode(requestHexByte[3:])
	requestHexByte = append(requestHexByte, verifyCode)
	log.WithField("请求内容", requestHexByte).Infof("开关设备 请求内容")

	_, err = conn.Write(requestHexByte)
	if err != nil {
		log.Warnf("开关设备 发送数据出错 %s", err.Error())
		return nil, err
	}
	reader := bufio.NewReader(conn)
	var buf [4096]byte
	n, err := reader.Read(buf[:]) // 读取数据
	if err != nil {
		log.Warnf("开关设备 接收返回消息出错 %s", err.Error())
		return nil, err
	}

	log.WithField("接收数据", buf[:n]).Info("开关设备 接收到数据")
	resp, err := getRespMsg(buf[:n])
	if err != nil {
		log.Warnf("开关设备 解析返回数据格式出错 %s", err.Error())
		return nil, err
	}
	log.WithField("接受数据", resp).Info("开关设备 返回成功")
	return resp, nil
}

func (agent *Agent) sendMsg(content string) (*string, error) {
	agent.mutex.Lock()
	defer agent.mutex.Unlock()
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

func getDeviceIdHex(deviceId string) (*string, error) {
	deviceIdInt, err := strconv.Atoi(deviceId)
	if err != nil {
		log.Warnf("转换设备号为十六进制字符串失败 %s", deviceId)
		return nil, err
	}
	hexStr := strconv.FormatInt(int64(deviceIdInt), 16)
	if len(hexStr) > 8 {
		log.Warnf("转换设备号为十六进制字符串 字符串过长 %s", deviceIdInt)
		return nil, errors.New("转换设备号为十六进制字符串 字符串过长")
	}
	for len(hexStr) < 8 {
		hexStr = "0" + hexStr
	}
	return &hexStr, nil
}

func byteToHexString(content []byte) string {
	result := ""
	for _, c := range content {
		hexStr := strconv.FormatInt(int64(c), 16)
		result = result + hexStr
	}
	return result
}

type ByteResp struct {
	Type           string `json:"type"`
	DeviceId       string `json:"deviceId"`
	Length         int    `json:"length"`
	DeviceInfoResp `json:"deviceInfoResp"`
}

type DeviceInfoResp struct {
	Signal     int `json:"signal"`     // 信号强度
	OpenStatus int `json:"openStatus"` // 开关状态
	OuterOpen  int `json:"outerOpen"`  // 外部开关量
	Light      int `json:"light"`      // 光照值
}

func getRespMsg(content []byte) (*ByteResp, error) {
	if len(content) < 11 {
		return nil, errors.New("解析包 返回包的长度不够")
	}
	if content[0] != byte(0x5a) || content[1] != byte(0xa5) {
		return nil, errors.New("解析包 包头部不正确")
	}
	length := int(content[2])
	if length != len(content)-4 {
		log.Warnf("解析包 返回包长度不等于实际包长度 包长度 %d 实际长度 %d", int(length), len(content)-3)
		return nil, errors.New("解析包 返回包长度不等于实际包长度")
	}
	result := &ByteResp{
		Length:   length,
		Type:     getCommand(content[8]),
		DeviceId: hex.EncodeToString(content[3:7]),
	}
	if result.Type == "" {
		return nil, errors.New("解析包 返回未知的操作指令")
	}
	if content[8] != COMMAND_GET_DEVICE_ID {
		result.Signal = int(content[9])
		result.OpenStatus = int(content[10])
		result.OuterOpen = int(content[11])
		result.Light = int(content[12])
	}
	return result, nil
}

func getCommand(command byte) string {
	if command == COMMAND_GET_DEVICE_ID {
		return COMMAND_GET_DEVICE_ID_STRING
	} else if command == COMMAND_GET_STATUS {
		return COMMAND_GET_STATUS_STRING
	} else if command == COMMAND_TURN_ON_OFF {
		return COMMAND_TURN_ON_OFF_STRING
	} else {
		return ""
	}
}

func genVerifyCode(verifyContent []byte) byte {
	result := byte(0)
	for _, c := range verifyContent {
		result = result + c
	}
	return result
}
