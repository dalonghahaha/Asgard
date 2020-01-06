package models

import "github.com/dalonghahaha/avenger/components/db"

type App struct {
	BaseModel
	OperatorModel
	CmdModel
	GroupID     int64  `gorm:"column:group_id" json:"group_id"`
	Name        string `gorm:"column:name" json:"name"`
	AutoRestart int64  `gorm:"column:auto_restart" json:"auto_restart"`
	IsMonitor   int64  `gorm:"column:is_monitor" json:"is_monitor"`
	Status      int64  `gorm:"column:status" json:"status"`
}

func (m *App) TableName() string {
	return "apps"
}

func (c *App) All() (list []*App, err error) {
	err = db.Get(DB_NAME).Find(&list).Error
	return
}

func (c *App) Search(where map[string]interface{}) (list []*App, err error) {
	err = db.Get(DB_NAME).Where(where).Find(&list).Error
	return
}

func (c *App) Get(id int64) (err error) {
	err = db.Get(DB_NAME).Where("id = ? ", id).First(c).Error
	return
}

func (c *App) Find(where map[string]interface{}) (err error) {
	err = db.Get(DB_NAME).Where(where).First(c).Error
	return
}

func (c *App) Create() (err error) {
	err = db.Get(DB_NAME).Create(c).Error
	return
}

func (c *App) Update() (err error) {
	err = db.Get(DB_NAME).Save(c).Error
	return
}

func (c *App) Delete() (err error) {
	err = db.Get(DB_NAME).Delete(c).Error
	return
}
