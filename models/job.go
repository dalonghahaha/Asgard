package models

import "github.com/dalonghahaha/avenger/components/db"

type Job struct {
	BaseModel
	OperatorModel
	CmdModel
	GroupID   int64  `gorm:"column:group_id" json:"group_id"`
	Name      string `gorm:"column:name" json:"name"`
	Spec      string `gorm:"column:spec" json:"spec"`
	Timeout   int64  `gorm:"column:timeout" json:"timeout"`
	IsMonitor int64  `gorm:"column:is_monitor" json:"is_monitor"`
	Status    int64  `gorm:"column:status" json:"status"`
}

func (m *Job) TableName() string {
	return "jobs"
}

func (c *Job) All() (list []*Job, err error) {
	err = db.Get(DB_NAME).Find(&list).Error
	return
}

func (c *Job) Search(where map[string]interface{}) (list []*Job, err error) {
	err = db.Get(DB_NAME).Where(where).Find(&list).Error
	return
}

func (c *Job) Get(id int64) (err error) {
	err = db.Get(DB_NAME).Where("id = ? ", id).First(c).Error
	return
}

func (c *Job) Find(where map[string]interface{}) (err error) {
	err = db.Get(DB_NAME).Where(where).First(c).Error
	return
}

func (c *Job) Create() (err error) {
	err = db.Get(DB_NAME).Create(c).Error
	return
}

func (c *Job) Update() (err error) {
	err = db.Get(DB_NAME).Save(c).Error
	return
}

func (c *Job) Delete() (err error) {
	err = db.Get(DB_NAME).Delete(c).Error
	return
}
