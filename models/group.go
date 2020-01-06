package models

import "github.com/dalonghahaha/avenger/components/db"

type Group struct {
	BaseModel
	OperatorModel
	Name   string `gorm:"column:name" json:"name"`
	Status int64  `gorm:"column:status" json:"status"`
}

func (m *Group) TableName() string {
	return "groups"
}

func (c *Group) All() (list []*Group, err error) {
	err = db.Get(DB_NAME).Find(&list).Error
	return
}

func (c *Group) Search(where map[string]interface{}) (list []*Group, err error) {
	err = db.Get(DB_NAME).Where(where).Find(&list).Error
	return
}

func (c *Group) Get(id int64) (err error) {
	err = db.Get(DB_NAME).Where("id = ? ", id).First(c).Error
	return
}

func (c *Group) Find(where map[string]interface{}) (err error) {
	err = db.Get(DB_NAME).Where(where).First(c).Error
	return
}

func (c *Group) Create() (err error) {
	err = db.Get(DB_NAME).Create(c).Error
	return
}

func (c *Group) Update() (err error) {
	err = db.Get(DB_NAME).Save(c).Error
	return
}

func (c *Group) Delete() (err error) {
	err = db.Get(DB_NAME).Delete(c).Error
	return
}
