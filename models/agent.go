package models

type Agent struct {
	BaseModel
	Alias  string `gorm:"column:alias" json:"alias"`
	IP     string `gorm:"column:ip" json:"ip"`
	Port   string `gorm:"column:port" json:"port"`
	Master string `gorm:"column:master" json:"master"`
	Status int64  `gorm:"column:status" json:"status"`
}

func (m *Agent) TableName() string {
	return "agents"
}
