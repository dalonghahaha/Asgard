package models

import "github.com/dalonghahaha/avenger/components/db"

type User struct {
	ID       int64  `gorm:"column:id;primary_key" json:"id"`
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

func (c *User) All() (list []*User, err error) {
	err = db.Get(DB_NAME).Find(&list).Error
	return
}

func (c *User) Search(where map[string]interface{}) (list []*User, err error) {
	err = db.Get(DB_NAME).Where(where).Find(&list).Error
	return
}

func (c *User) Get(id int64) (err error) {
	err = db.Get(DB_NAME).Where("id = ? ", id).First(c).Error
	return
}

func (c *User) Find(where map[string]interface{}) (err error) {
	err = db.Get(DB_NAME).Where(where).First(c).Error
	return
}

func (c *User) Create() (err error) {
	err = db.Get(DB_NAME).Create(c).Error
	return
}

func (c *User) Update() (err error) {
	err = db.Get(DB_NAME).Save(c).Error
	return
}

func (c *User) Delete() (err error) {
	err = db.Get(DB_NAME).Delete(c).Error
	return
}
