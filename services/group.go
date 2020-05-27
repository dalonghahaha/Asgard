package services

import (
	"encoding/json"
	"fmt"

	"github.com/dalonghahaha/avenger/components/logger"
	"github.com/jinzhu/gorm"

	"Asgard/constants"
	"Asgard/models"
)

type GroupService struct {
}

func NewGroupService() *GroupService {
	return &GroupService{}
}

func (s *GroupService) buidCacheKey(id int64) string {
	return fmt.Sprintf("%s:%d", constants.CACHE_KEY_GROUP, id)
}

func (s *GroupService) GetUsageGroup() (list []models.Group) {
	err := models.Where(&list, "status = ?", "1")
	if err != nil {
		logger.Error("GetUsageGroup Error:", err)
		return nil
	}
	return
}

func (s *GroupService) GetGroupPageList(where map[string]interface{}, page int, pageSize int) (list []models.Group, count int) {
	condition := "1=1"
	for key, val := range where {
		if key == "status" {
			if val.(int) == -99 {
				condition += " and status != -1"
			} else {
				condition += fmt.Sprintf(" and %s=%v", key, val)
			}
		} else if key == "name" {
			condition += fmt.Sprintf(" and %s like '%%%v%%' ", key, val)
		} else {
			condition += fmt.Sprintf(" and %s=%v", key, val)
		}
	}
	err := models.PageListbyWhereString(&models.Group{}, condition, page, pageSize, "created_at desc", &list, &count)
	if err != nil {
		logger.Error("GetGroupPageList Error:", err)
		return nil, 0
	}
	return
}

func (s *GroupService) GetGroupByID(id int64) *models.Group {
	var group models.Group
	data := GetCache(s.buidCacheKey(id))
	if len(data) > 0 {
		err := json.Unmarshal([]byte(data), &group)
		if err == nil {
			return &group
		}
	}
	err := models.Get(id, &group)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			logger.Error("GetUserByID Error:", err)
		}
		return nil
	}
	info, _ := json.Marshal(group)
	SetCache(s.buidCacheKey(id), string(info))
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
	DelCache(s.buidCacheKey(group.ID))
	return true
}

func (s *GroupService) DeleteGroup(group *models.Group) bool {
	err := models.Delete(group)
	if err != nil {
		logger.Error("DeleteGroupByID Error:", err)
		return false
	}
	DelCache(s.buidCacheKey(group.ID))
	return true
}

func (s *GroupService) ChangeGroupStatus(group *models.Group, status int64, updator int64) bool {
	values := map[string]interface{}{
		"status":  status,
		"updator": updator,
	}
	err := models.UpdateColumns(group, values)
	if err != nil {
		logger.Error("ChangeGroupStatus Error:", err)
		return false
	}
	return true
}
