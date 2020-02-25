package models

import (
	"time"
)

type Timing struct {
	BaseModel
	OperatorModel
	CmdModel
	GroupID   int64     `gorm:"column:group_id" json:"group_id"`
	Name      string    `gorm:"column:name" json:"name"`
	AgentID   int64     `gorm:"column:agent_id" json:"agent_id"`
	Time      time.Time `gorm:"column:time" json:"time"`
	Timeout   int64     `gorm:"column:timeout" json:"timeout"`
	IsMonitor int64     `gorm:"column:is_monitor" json:"is_monitor"`
	Status    int64     `gorm:"column:status" json:"status"`
}

func (m *Timing) TableName() string {
	return "timings"
}
