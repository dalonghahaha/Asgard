package services

import (
	"encoding/json"
	"fmt"

	"github.com/dalonghahaha/avenger/components/logger"
	"github.com/jinzhu/gorm"

	"Asgard/constants"
	"Asgard/models"
)

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) buidCacheKey(id int64) string {
	return fmt.Sprintf("%s:%d", constants.CACHE_KEY_USER, id)
}

func (s *UserService) GetUserPageList(where map[string]interface{}, page int, pageSize int) (list []models.User, count int) {
	condition := "1=1"
	for key, val := range where {
		if key == "status" {
			if val.(int) == -99 {
				condition += " and status != -1"
			} else {
				condition += fmt.Sprintf(" and %s=%v", key, val)
			}
		} else if key == "nickname" || key == "phone" || key == "email" {
			condition += fmt.Sprintf(" and %s like '%%%v%%' ", key, val)
		} else {
			condition += fmt.Sprintf(" and %s=%v", key, val)
		}
	}
	err := models.PageListbyWhereString(&models.User{}, condition, page, pageSize, "created_at desc", &list, &count)
	if err != nil {
		logger.Error("GetUserPageList Error:", err)
		return nil, 0
	}
	return
}

func (s *UserService) GetUserByID(id int64) *models.User {
	var user models.User
	data := GetCache(s.buidCacheKey(id))
	if len(data) > 0 {
		err := json.Unmarshal([]byte(data), &user)
		if err == nil {
			return &user
		}
	}
	err := models.Get(id, &user)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			logger.Error("GetUserByID Error:", err)
		}
		return nil
	}
	info, _ := json.Marshal(user)
	SetCache(s.buidCacheKey(id), string(info))
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
	DelCache(s.buidCacheKey(user.ID))
	return true
}

func (s *UserService) DeleteUserByID(user *models.User) bool {
	err := models.Delete(user)
	if err != nil {
		logger.Error("DeleteUserByID Error:", err)
		return false
	}
	DelCache(s.buidCacheKey(user.ID))
	return true
}
