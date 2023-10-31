package entity

import (
	"gorm.io/gorm"
	"time"
)

type Device struct {
	Id         int64     `gorm:"primaryKey" json:"id"`
	Name       string    `json:"name"`
	DeviceId   string    `json:"deviceId"`
	Location   string    `json:"location"`
	Status     string    `json:"status"`
	Type       string    `json:"type"`
	CreateBy   string    `json:"createBy"`
	CreateTime time.Time `json:"createTime" gorm:"type：timestamp"`
	UpdateBy   string    `json:"updateBy"`
	UpdateTime time.Time `json:"updateTime" gorm:"type：timestamp"`
}

func (Device) TableName() string {
	return deviceTableName
}

const (
	deviceTableName = "device"
	updateBy        = "数据定时上报"
	createBy        = "数据上报"
	LAMP_TYPE       = "智能照明设备"
)

func ChangeDeviceState(id int64, status string) *gorm.DB {
	return db.Table(deviceTableName).Where("id = ?", id).UpdateColumns(Device{
		Status:     status,
		UpdateBy:   updateBy,
		UpdateTime: time.Now(),
	})
}

func GetByDeviceId(deviceId string) (*Device, error) {
	device := &Device{}
	tx := db.Table(deviceTableName).Where("device_id = ?", deviceId).First(device)
	return device, tx.Error
}

func CreateDevice(device *Device) *gorm.DB {
	if device == nil {
		return nil
	}
	device.CreateBy = createBy
	device.CreateTime = time.Now()
	return db.Table(deviceTableName).Create(device)
}

func BatchCreateDevice(devices []*Device) *gorm.DB {
	if devices == nil || len(devices) == 0 {
		return nil
	}
	for _, device := range devices {
		device.CreateBy = createBy
		device.CreateTime = time.Now()
	}
	return db.Table(deviceTableName).Create(devices)
}

func ListDeviceByCursor(id int64, limit int) []*Device {
	devices := make([]*Device, limit)
	db.Table(deviceTableName).Select("*").Where("id > ?", id).Limit(limit).Find(&devices)
	return devices
}

func ListLampDeviceByCursor(id int64, limit int) []*Device {
	devices := make([]*Device, limit)
	db.Table(deviceTableName).Select("*").Where("id > ?", id).Where("type = ?", LAMP_TYPE).Limit(limit).Find(&devices)
	return devices
}
