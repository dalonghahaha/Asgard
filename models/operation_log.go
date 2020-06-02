package models

import "time"

type OperationLog struct {
	ID        int64     `gorm:"column:id;primary_key" json:"id"`
	UserID    int64     `gorm:"column:user_id" json:"user_id"`
	Type      int64     `gorm:"column:type" json:"type"`
	RelatedID int64     `gorm:"column:related_id" json:"related_id"`
	Action    int64     `gorm:"column:action" json:"action"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

func (m *OperationLog) TableName() string {
	return "operation_logs"
}
