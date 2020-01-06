package models

import "github.com/dalonghahaha/avenger/components/db"

type AgentJob struct {
	ID      int64 `gorm:"column:id;primary_key" json:"id"`
	AgentID int64 `gorm:"column:agent_id" json:"agent_id"`
	JobID   int64 `gorm:"column:job_id" json:"job_id"`
}

func (m *AgentJob) TableName() string {
	return "agent_jobs"
}

func (c *AgentJob) All() (list []*AgentJob, err error) {
	err = db.Get(DB_NAME).Find(&list).Error
	return
}

func (c *AgentJob) Search(where map[string]interface{}) (list []*AgentJob, err error) {
	err = db.Get(DB_NAME).Where(where).Find(&list).Error
	return
}

func (c *AgentJob) Get(id int64) (err error) {
	err = db.Get(DB_NAME).Where("id = ? ", id).First(c).Error
	return
}

func (c *AgentJob) Find(where map[string]interface{}) (err error) {
	err = db.Get(DB_NAME).Where(where).First(c).Error
	return
}

func (c *AgentJob) Create() (err error) {
	err = db.Get(DB_NAME).Create(c).Error
	return
}

func (c *AgentJob) Update() (err error) {
	err = db.Get(DB_NAME).Save(c).Error
	return
}

func (c *AgentJob) Delete() (err error) {
	err = db.Get(DB_NAME).Delete(c).Error
	return
}
