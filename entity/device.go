package entity

import (
	"gorm.io/gorm"
)

type Device struct {
	Id       int64  `gorm:"primaryKey" json:"id"`
	Name     string `json:"name"`
	DeviceId string `json:"deviceId"`
	Location string `json:"location"`
	Status   string `json:"status"`
}

func (Device) TableName() string {
	return deviceTableName
}

const deviceTableName = "device"

func ChangeDeviceState(id int64, status string) *gorm.DB {
	return db.Table(deviceTableName).Where("id = ?", id).UpdateColumn("status", status)
}

func GetByDeviceId(deviceId string) *Device {
	var device *Device
	db.Table(deviceTableName).Where("device_id = ?", deviceId).First(device)
	return device
}

func CreateDevice(device *Device) *gorm.DB {
	return db.Table(deviceTableName).Create(device)
}

func BatchCreateDevice(devices []*Device) *gorm.DB {
	return db.Table(deviceTableName).Create(devices)
}

func ListDeviceByCursor(id int64, limit int) []*Device {
	devices := make([]*Device, limit)
	db.Table(deviceTableName).Select("*").Where("id > ?", id).Limit(limit).Find(&devices)
	return devices
}
