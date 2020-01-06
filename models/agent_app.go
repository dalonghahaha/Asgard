package models

import "github.com/dalonghahaha/avenger/components/db"

type AgentApp struct {
	ID      int64 `gorm:"column:id;primary_key" json:"id"`
	AgentID int64 `gorm:"column:agent_id" json:"agent_id"`
	AppID   int64 `gorm:"column:app_id" json:"app_id"`
}

func (m *AgentApp) TableName() string {
	return "agent_apps"
}

func (c *AgentApp) All() (list []*AgentApp, err error) {
	err = db.Get(DB_NAME).Find(&list).Error
	return
}

func (c *AgentApp) Search(where map[string]interface{}) (list []*AgentApp, err error) {
	err = db.Get(DB_NAME).Where(where).Find(&list).Error
	return
}

func (c *AgentApp) Get(id int64) (err error) {
	err = db.Get(DB_NAME).Where("id = ? ", id).First(c).Error
	return
}

func (c *AgentApp) Find(where map[string]interface{}) (err error) {
	err = db.Get(DB_NAME).Where(where).First(c).Error
	return
}

func (c *AgentApp) Create() (err error) {
	err = db.Get(DB_NAME).Create(c).Error
	return
}

func (c *AgentApp) Update() (err error) {
	err = db.Get(DB_NAME).Save(c).Error
	return
}

func (c *AgentApp) Delete() (err error) {
	err = db.Get(DB_NAME).Delete(c).Error
	return
}
