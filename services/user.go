package services

import (
	"github.com/dalonghahaha/avenger/components/logger"
	"github.com/jinzhu/gorm"

	"Asgard/models"
)

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) GetAllUser() []*models.User {
	list, err := new(models.User).All()
	if err != nil {
		logger.Error("GetAllUser Error:", err)
		return nil
	}
	return list
}

func (s *UserService) GetUserByID(id int64) *models.User {
	user := new(models.User)
	err := user.Get(id)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			logger.Error("GetUserByID Error:", err)
		}
		return nil
	}
	return user
}

func (s *UserService) GetUserByNickName(nickname string) *models.User {
	where := map[string]interface{}{
		"nickname": nickname,
	}
	user := new(models.User)
	err := user.Find(where)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			logger.Error("GetUserByNickName Error:", err)
		}
		return nil
	}
	return user
}

func (s *UserService) GetUserByEmail(email string) *models.User {
	where := map[string]interface{}{
		"email": email,
	}
	user := new(models.User)
	err := user.Find(where)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			logger.Error("GetUserByEmail Error:", err)
		}
		return nil
	}
	return user
}

func (s *UserService) GetUserByMobile(mobile string) *models.User {
	where := map[string]interface{}{
		"mobile": mobile,
	}
	user := new(models.User)
	err := user.Find(where)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			logger.Error("GetUserByMobile Error:", err)
		}
		return nil
	}
	return user
}

func (s *UserService) CreateUser(user *models.User) bool {
	err := user.Create()
	if err != nil {
		logger.Error("CreateUser Error:", err)
		return false
	}
	return true
}

func (s *UserService) UpdateUser(user *models.User) bool {
	err := user.Update()
	if err != nil {
		logger.Error("UpdateUser Error:", err)
		return false
	}
	return true
}

func (s *UserService) DeleteUserByID(id int64) bool {
	user := new(models.User)
	user.ID = id
	err := user.Delete()
	if err != nil {
		logger.Error("DeleteUserByID Error:", err)
		return false
	}
	return true
}
