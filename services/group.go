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

func (s *GroupService) GetAllGroup() (list []models.Group) {
	err := models.All(&list)
	if err != nil {
		logger.Error("GetAllGroup Error:", err)
		return nil
	}
	return
}

func (s *GroupService) GetGroupPageList(where map[string]interface{}, page int, pageSize int) (list []models.Group, count int) {
	err := models.PageList(&models.Group{}, where, page, pageSize, &list, &count)
	if err != nil {
		logger.Error("GetGroupPageList Error:", err)
		return nil, 0
	}
	return
}

func (s *GroupService) GetGroupByID(id int64) *models.Group {
	var group models.Group
	err := models.Get(id, &group)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			logger.Error("GetUserByID Error:", err)
		}
		return nil
	}
	return &group
}

func (s *GroupService) CreateGroup(group *models.Group) bool {
	err := models.Create(group)
	if err != nil {
		logger.Error("CreateGroup Error:", err)
		return false
	}
	return true
}

func (s *GroupService) UpdateGroup(group *models.Group) bool {
	err := models.Update(group)
	if err != nil {
		logger.Error("UpdateGroup Error:", err)
		return false
	}
	return true
}

func (s *GroupService) DeleteGroupByID(id int64) bool {
	group := new(models.Group)
	group.ID = id
	err := models.Delete(group)
	if err != nil {
		logger.Error("DeleteGroupByID Error:", err)
		return false
	}
	return true
}
