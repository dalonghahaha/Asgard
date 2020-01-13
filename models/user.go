package models

type User struct {
	BaseModel
	NickName string `gorm:"column:nickname" json:"nickname"`
	Avatar   string `gorm:"column:avatar" json:"avatar"`
	Email    string `gorm:"column:email" json:"email"`
	Mobile   string `gorm:"column:mobile" json:"mobile"`
	Salt     string `gorm:"column:salt" json:"salt"`
	Password string `gorm:"column:password" json:"password"`
	Status   int64  `gorm:"column:status" json:"status"`
}

func (m *User) TableName() string {
	return "users"
}
