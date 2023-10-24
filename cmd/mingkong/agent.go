package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
)

var deviceIdAgentPool map[string]*MingAgent

type MingAgent struct {
	DeviceId string
	Status   string
}

type MingReq struct {
	Command  string
	DeviceId string
	Time     int64
}

type MingResp struct {
}

func byteServe() {
	listen, err := net.Listen("tcp4", "0.0.0.0:8005")
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

		reader := bufio.NewReader(conn)
		var buf [4096]byte
		n, err := reader.Read(buf[:]) // 读取数据
		if err != nil {
			log.Warnf("读取请求体出错 %s", err.Error())
			_ = conn.Close()
			continue
		}
		handle(buf[:n])
	}
}

func handle(reqContent []byte) {
	log.WithField("接收请求数据", reqContent).Info("接收请求数据")
	if reqContent[0] != byte(0xA5) || reqContent[1] != byte(0x5A) || reqContent[2] != byte(0x04) {
		log.Warnf("解析包 包头部不正确")
		return
	}
	if reqContent[len(reqContent)-1] != byte(0xAA) || reqContent[len(reqContent)-2] != byte(0x55) {
		log.Warnf("解析包 包结尾不正确")
		return
	}
	length := int(reqContent[3])*256 + int(reqContent[4])
	if length != len(reqContent) {
		log.Warnf("解析包 返回包长度不等于实际包长度 包长度 %d 实际长度 %d", length, len(reqContent))
		return
	}
	deviceId := hex.EncodeToString(reqContent[5:18])
	dataLength := int(reqContent[21])*256 + int(reqContent[22])
	if dataLength == 0 {
		log.Warnf("上报数据长度为0")
		return
	}
	data := reqContent[23 : 23+dataLength]
	for i := 0; i < dataLength; i++ {

	}
}

func refreshOnline() {
	
}
