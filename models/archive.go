package models

import (
	"time"
)

type Archive struct {
	ID        int64     `gorm:"column:id;primary_key" json:"id"`
	Type      int64     `gorm:"column:type" json:"type"`
	RelatedID int64     `gorm:"column:related_id" json:"related_id"`
	UUID      string    `gorm:"column:uuid" json:"uuid"`
	PID       int64     `gorm:"column:pid" json:"pid"`
	BeginTime time.Time `gorm:"column:begin_time" json:"begin_time"`
	EndTime   time.Time `gorm:"column:end_time" json:"end_time"`
	Status    int64     `gorm:"column:status" json:"status"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

func (m *Archive) TableName() string {
	return "archives"
}
