package models

import "time"

type Exception struct {
	ID        int64     `gorm:"column:id;primary_key" json:"id"`
	Type      int64     `gorm:"column:type" json:"type"`
	RelatedID int64     `gorm:"column:related_id" json:"related_id"`
	Desc      string    `gorm:"column:desc" json:"desc"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

func (m *Exception) TableName() string {
	return "exceptions"
}
