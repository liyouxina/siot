package main

import (
	"github.com/liyouxina/siot/entity"
	log "github.com/sirupsen/logrus"
	"time"
)

func monitor() {
	for {
		time.Sleep(time.Minute)
		scanAndKeepDevices()
	}
}

func scanAndKeepDevices() {
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
	agent := deviceIdAgentPool[device.DeviceId]
	if agent == nil {
		tx := entity.ChangeDeviceState(device.Id, "没找到这个设备连接")
		if tx.Error != nil {
			log.Warnf("设备心跳检测 改变设备状态为没找到这个设备连接 返回错误 %s", tx.Error.Error())
		}
		if tx.RowsAffected == 0 {
			log.Warnf("设备心跳检测 改变设备状态为没找到这个设备连接 失败")
		}
		return
	}
	// 心跳检测
	resp, err := agent.GetDeviceId()
	if err != nil {
		log.Warnf("设备心跳检测 返回有错误 %s", err.Error())
		tx := entity.ChangeDeviceState(device.Id, "设备心跳检测有问题")
		if tx.Error != nil {
			log.Warnf("设备心跳检测 返回有错误 改变设备状态为断开 返回错误 %s", tx.Error.Error())
		}
		if tx.RowsAffected == 0 {
			log.Warnf("设备心跳检测 返回有错误 改变设备状态为断开 失败")
		}
		agent.mutex.Lock()
		deviceIdAgentPool[device.DeviceId] = nil
		systemIdAgentPool[device.DeviceId] = nil
		agent.mutex.Unlock()
	}

	if resp.DeviceId != agent.DeviceId {
		log.Warnf("设备心跳检测 返回有错误 %s", err.Error())
		tx := entity.ChangeDeviceState(device.Id, "设备心跳检测有问题")
		if tx.Error != nil {
			log.Warnf("设备心跳检测 返回有错误 改变设备状态为断开 返回错误 %s", tx.Error.Error())
		}
		if tx.RowsAffected == 0 {
			log.Warnf("设备心跳检测 返回有错误 改变设备状态为断开 失败")
		}
		agent.mutex.Lock()
		deviceIdAgentPool[device.DeviceId] = nil
		systemIdAgentPool[device.DeviceId] = nil
		agent.mutex.Unlock()
	}
	log.Infof("设备心跳检测成功 设备号 %s", device.DeviceId)
}
