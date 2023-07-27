package entity

import "time"

type Device struct {
	id            int64
	name          string
	extra         string
	isHumanRecord bool
	createdAt     time.Time
	updatedAt     time.Time
	deletedAt     time.Time
	operator      string
}

func (*Device) List() {

}