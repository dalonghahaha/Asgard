package models

import "github.com/dalonghahaha/avenger/components/db"

type Agent struct {
	BaseModel
	IP     string `gorm:"column:ip" json:"ip"`
	Port   string `gorm:"column:port" json:"port"`
	Status int64  `gorm:"column:status" json:"status"`
}

func (m *Agent) TableName() string {
	return "agents"
}

func (c *Agent) All() (list []*Agent, err error) {
	err = db.Get(DB_NAME).Find(&list).Error
	return
}

func (c *Agent) Search(where map[string]interface{}) (list []*Agent, err error) {
	err = db.Get(DB_NAME).Where(where).Find(&list).Error
	return
}

func (c *Agent) Get(id int64) (err error) {
	err = db.Get(DB_NAME).Where("id = ? ", id).First(c).Error
	return
}

func (c *Agent) Find(where map[string]interface{}) (err error) {
	err = db.Get(DB_NAME).Where(where).First(c).Error
	return
}

func (c *Agent) Create() (err error) {
	err = db.Get(DB_NAME).Create(c).Error
	return
}

func (c *Agent) Update() (err error) {
	err = db.Get(DB_NAME).Save(c).Error
	return
}

func (c *Agent) Delete() (err error) {
	err = db.Get(DB_NAME).Delete(c).Error
	return
}
