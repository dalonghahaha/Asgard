package models

import "time"

type Monitor struct {
	ID        int64     `gorm:"column:id;primary_key" json:"id"`
	Type      int64     `gorm:"column:type" json:"type"`
	RelatedID int64     `gorm:"column:related_id" json:"related_id"`
	UUID      string    `gorm:"column:uuid" json:"uuid"`
	PID       int64     `gorm:"column:pid" json:"pid"`
	CPU       float64   `gorm:"column:cpu" json:"cpu"`
	Memory    float64   `gorm:"column:memory" json:"memory"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

func (m *Monitor) TableName() string {
	return "monitors"
}
