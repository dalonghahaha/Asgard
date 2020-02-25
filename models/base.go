package models

import (
	"time"
)

const (
	DB_NAME     = "asgard"
	TYPE_AGENT  = int64(1)
	TYPE_APP    = int64(2)
	TYPE_JOB    = int64(3)
	TYPE_TIMING = int64(4)
)

//Model基类，定义通用属性
type BaseModel struct {
	ID        int64     `gorm:"column:id;primary_key" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type OperatorModel struct {
	Creator int64 `gorm:"column:creator" json:"creator"`
	Updator int64 `gorm:"column:updator" json:"updator"`
}

type CmdModel struct {
	Dir     string `gorm:"column:dir" json:"dir"`
	Program string `gorm:"column:program" json:"program"`
	Args    string `gorm:"column:args" json:"args"`
	StdOut  string `gorm:"column:std_out" json:"std_out"`
	StdErr  string `gorm:"column:std_err" json:"std_err"`
}

//BaseModel的BeforeCreate钩子
func (c *BaseModel) BeforeCreate() error {
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	return nil
}

//BaseModel的BeforeUpdate钩子
func (c *BaseModel) BeforeUpdate() error {
	c.UpdatedAt = time.Now()
	return nil
}
