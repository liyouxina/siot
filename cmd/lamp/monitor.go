package main

import (
	"github.com/liyouxina/siot/entity"
	log "github.com/sirupsen/logrus"
	"time"
)

func monitor() {
	for {
		time.Sleep(time.Second * 10)
		scanAndKeepDevices()
	}
}

func scanAndKeepDevices() {
	cursor := int64(0)
	for {
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
	log.Infof("设备心跳检测 开始检测设备号 %s", device.DeviceId)
	agent := deviceIdAgentPool[device.DeviceId]
	if agent == nil {
		tx := entity.ChangeDeviceState(device.Id, "没找到这个设备连接")
		if tx.Error != nil {
			log.Warnf("设备心跳检测 改变设备状态为没找到这个设备连接 返回错误 %s %s", device.DeviceId, tx.Error.Error())
		}
		if tx.RowsAffected == 0 {
			log.Warnf("设备心跳检测 改变设备状态为没找到这个设备连接 失败 %s", device.DeviceId)
		}
		return
	}
	// 心跳检测
	resp, err := agent.GetDeviceId()
	if err != nil {
		log.Warnf("设备心跳检测 返回有错误 %s %s", device.DeviceId, err.Error())
		tx := entity.ChangeDeviceState(device.Id, "设备心跳检测有问题")
		if tx.Error != nil {
			log.Warnf("设备心跳检测 返回有错误 改变设备状态为断开 返回错误 %s %s", device.DeviceId, tx.Error.Error())
		}
		if tx.RowsAffected == 0 {
			log.Warnf("设备心跳检测 返回有错误 改变设备状态为断开 失败 %s", device.DeviceId)
		}
		agent.mutex.Lock()
		deviceIdAgentPool[device.DeviceId] = nil
		systemIdAgentPool[device.DeviceId] = nil
		agent.mutex.Unlock()
		return
	}
	if resp == nil {
		log.Warnf("设备心跳检测 返回值为空 %s", device.DeviceId)
		agent.mutex.Lock()
		deviceIdAgentPool[device.DeviceId] = nil
		systemIdAgentPool[device.DeviceId] = nil
		agent.mutex.Unlock()
		return
	}

	if resp.DeviceId != agent.DeviceId {
		log.Warnf("设备心跳检测 返回有错误 %s %s", device.DeviceId, err.Error())
		tx := entity.ChangeDeviceState(device.Id, "设备心跳检测有问题")
		if tx.Error != nil {
			log.Warnf("设备心跳检测 返回有错误 改变设备状态为断开 返回错误 %s %s", device.DeviceId, tx.Error.Error())
		}
		if tx.RowsAffected == 0 {
			log.Warnf("设备心跳检测 返回有错误 改变设备状态为断开 失败 %s", device.DeviceId)
		}
		agent.mutex.Lock()
		deviceIdAgentPool[device.DeviceId] = nil
		systemIdAgentPool[device.DeviceId] = nil
		agent.mutex.Unlock()
		return
	}
	log.Infof("设备心跳检测成功 设备号 %s", device.DeviceId)
}
