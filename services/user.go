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

func (s *UserService) GetUserPageList(where map[string]interface{}, page int, pageSize int) (list []models.User, count int) {
	err := models.PageList(&models.User{}, where, page, pageSize, "created_at desc", &list, &count)
	if err != nil {
		logger.Error("GetJobPageList Error:", err)
		return nil, 0
	}
	return
}

func (s *UserService) GetUserByID(id int64) *models.User {
	var user models.User
	err := models.Get(id, &user)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			logger.Error("GetUserByID Error:", err)
		}
		return nil
	}
	return &user
}

func (s *UserService) GetUserByNickName(nickname string) *models.User {
	where := map[string]interface{}{
		"nickname": nickname,
	}
	var user models.User
	err := models.Find(where, &user)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			logger.Error("GetUserByNickName Error:", err)
		}
		return nil
	}
	return &user
}

func (s *UserService) GetUserByEmail(email string) *models.User {
	where := map[string]interface{}{
		"email": email,
	}
	var user models.User
	err := models.Find(where, &user)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			logger.Error("GetUserByEmail Error:", err)
		}
		return nil
	}
	return &user
}

func (s *UserService) GetUserByMobile(mobile string) *models.User {
	where := map[string]interface{}{
		"mobile": mobile,
	}
	var user models.User
	err := models.Find(where, &user)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			logger.Error("GetUserByMobile Error:", err)
		}
		return nil
	}
	return &user
}

func (s *UserService) CreateUser(user *models.User) bool {
	err := models.Create(user)
	if err != nil {
		logger.Error("CreateUser Error:", err)
		return false
	}
	return true
}

func (s *UserService) UpdateUser(user *models.User) bool {
	err := models.Update(user)
	if err != nil {
		logger.Error("UpdateUser Error:", err)
		return false
	}
	return true
}

func (s *UserService) DeleteUserByID(id int64) bool {
	user := new(models.User)
	user.ID = id
	err := models.Delete(user)
	if err != nil {
		logger.Error("DeleteUserByID Error:", err)
		return false
	}
	return true
}
