package models

import (
	"time"
)

type Archive struct {
	ID        int64     `gorm:"column:id;primary_key" json:"id"`
	Type      int64     `gorm:"column:type" json:"type"`
	RelatedID string    `gorm:"column:related_id" json:"related_id"`
	UUID      string    `gorm:"column:uuid" json:"uuid"`
	BeginTime time.Time `gorm:"column:begin_time" json:"begin_time"`
	EndTime   time.Time `gorm:"column:end_time" json:"end_time"`
	Status    int64     `gorm:"column:status" json:"status"`
}

func (m *Archive) TableName() string {
	return "archives"
}
