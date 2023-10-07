package main

import (
	"fmt"
	"github.com/liyouxina/siot/entity"
	log "github.com/sirupsen/logrus"
	"net"
	"strconv"
	"sync"
	"time"
)

var deviceIdAgentPool map[string]*MingAgent
var systemIdAgentPool map[string]*MingAgent

type MingAgent struct {
	DeviceId string
	SystemId string
	Status   string
	mutex    sync.Mutex
}

type MingReq struct {
	Command  string
	DeviceId string
	Time     int64
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
			log.Warnf("接收连接出错 %s", err.Error())
			continue
		}
		log.Infof("接收连接 接收到连接")

		agent := MingAgent{
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
		if err != nil {
			log.Warnf("接受连接 读取数据库出错 %s %s %s", agent.SystemId, agent.DeviceId, err.Error())
			continue
		}
		if deviceDO == nil {
			deviceDO = &entity.Device{
				Name:     "lamp " + agent.SystemId,
				DeviceId: agent.DeviceId,
				Status:   agent.Status,
			}
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
