package models

type Monitor struct {
	ID        int64   `gorm:"column:id;primary_key" json:"id"`
	Type      int64   `gorm:"column:type" json:"type"`
	RelatedID string  `gorm:"column:related_id" json:"related_id"`
	PID       string  `gorm:"column:pid" json:"pid"`
	CPU       float64 `gorm:"column:cpu" json:"cpu"`
	Memory    float64 `gorm:"column:memory" json:"memory"`
	Status    int64   `gorm:"column:status" json:"status"`
}

func (m *Monitor) TableName() string {
	return "monitors"
}
