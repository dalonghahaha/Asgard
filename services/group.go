package services

import (
	"github.com/dalonghahaha/avenger/components/logger"
	"github.com/jinzhu/gorm"

	"Asgard/models"
)

type GroupService struct {
}

func NewGroupService() *GroupService {
	return &GroupService{}
}

func (s *GroupService) GetAllGroup() []*models.Group {
	list, err := new(models.Group).All()
	if err != nil {
		logger.Error("GetAllGroup Error:", err)
		return nil
	}
	return list
}

func (s *GroupService) GetGroupByID(id int64) *models.Group {
	user := new(models.Group)
	err := user.Get(id)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			logger.Error("GetUserByID Error:", err)
		}
		return nil
	}
	return user
}

func (s *GroupService) CreateGroup(group *models.Group) bool {
	err := group.Create()
	if err != nil {
		logger.Error("CreateGroup Error:", err)
		return false
	}
	return true
}

func (s *GroupService) UpdateGroup(group *models.Group) bool {
	err := group.Update()
	if err != nil {
		logger.Error("UpdateGroup Error:", err)
		return false
	}
	return true
}

func (s *GroupService) DeleteGroupByID(id int64) bool {
	group := new(models.Group)
	group.ID = id
	err := group.Delete()
	if err != nil {
		logger.Error("DeleteGroupByID Error:", err)
		return false
	}
	return true
}
