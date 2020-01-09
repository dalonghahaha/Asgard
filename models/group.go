package models

type Group struct {
	BaseModel
	OperatorModel
	Name   string `gorm:"column:name" json:"name"`
	Status int64  `gorm:"column:status" json:"status"`
}

func (m *Group) TableName() string {
	return "groups"
}
