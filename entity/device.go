package entity

import (
	"gorm.io/gorm"
)

type Device struct {
	Id       int64  `gorm:"primaryKey" json:"id"`
	Name     string `json:"name"`
	DeviceId string `json:"deviceId"`
	Location string `json:"location"`
	State    string `json:"state"`
}

func (Device) TableName() string {
	return tableName
}

const tableName = "device"

func CreateDevice(device *Device) *gorm.DB {
	return db.Create(device)
}

func BatchCreateDevice(devices []*Device) *gorm.DB {
	return db.Create(devices)
}

func ListDeviceByCursor(id int64, limit int) []*Device {
	devices := make([]*Device, limit)
	db.Table(tableName).Select("*").Where("id > ?", id).Limit(limit).Find(&devices)
	return devices
}
