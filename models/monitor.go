package models

import (
	"github.com/dalonghahaha/avenger/components/db"
)

type Monitor struct {
	ID        int64   `gorm:"column:id;primary_key" json:"id"`
	Type      int64   `gorm:"column:type" json:"type"`
	RelatedID string  `gorm:"column:related_id" json:"related_id"`
	PID       string  `gorm:"column:pid" json:"pid"`
	CPU       float64 `gorm:"column:cpu" json:"cpu"`
	Memory    float64 `gorm:"column:memory" json:"memory"`
	Status    int64   `gorm:"column:status" json:"status"`
}

func (m *Monitor) TableName() string {
	return "monitors"
}

func (c *Monitor) All() (list []*Monitor, err error) {
	err = db.Get(DB_NAME).Find(&list).Error
	return
}

func (c *Monitor) Search(where map[string]interface{}) (list []*Monitor, err error) {
	err = db.Get(DB_NAME).Where(where).Find(&list).Error
	return
}

func (c *Monitor) Get(id int64) (err error) {
	err = db.Get(DB_NAME).Where("id = ? ", id).First(c).Error
	return
}

func (c *Monitor) Find(where map[string]interface{}) (err error) {
	err = db.Get(DB_NAME).Where(where).First(c).Error
	return
}

func (c *Monitor) Create() (err error) {
	err = db.Get(DB_NAME).Create(c).Error
	return
}

func (c *Monitor) Update() (err error) {
	err = db.Get(DB_NAME).Save(c).Error
	return
}

func (c *Monitor) Delete() (err error) {
	err = db.Get(DB_NAME).Delete(c).Error
	return
}
