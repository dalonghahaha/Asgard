package models

import (
	"time"

	"github.com/dalonghahaha/avenger/components/db"
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

func (c *Archive) All() (list []*Archive, err error) {
	err = db.Get(DB_NAME).Find(&list).Error
	return
}

func (c *Archive) Search(where map[string]interface{}) (list []*Archive, err error) {
	err = db.Get(DB_NAME).Where(where).Find(&list).Error
	return
}

func (c *Archive) Get(id int64) (err error) {
	err = db.Get(DB_NAME).Where("id = ? ", id).First(c).Error
	return
}

func (c *Archive) Find(where map[string]interface{}) (err error) {
	err = db.Get(DB_NAME).Where(where).First(c).Error
	return
}

func (c *Archive) Create() (err error) {
	err = db.Get(DB_NAME).Create(c).Error
	return
}

func (c *Archive) Update() (err error) {
	err = db.Get(DB_NAME).Save(c).Error
	return
}

func (c *Archive) Delete() (err error) {
	err = db.Get(DB_NAME).Delete(c).Error
	return
}
