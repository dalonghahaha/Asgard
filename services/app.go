package services

import (
	"fmt"

	"github.com/dalonghahaha/avenger/components/logger"
	"github.com/jinzhu/gorm"

	"Asgard/constants"
	"Asgard/models"
)

type AppService struct {
}

func NewAppService() *AppService {
	return &AppService{}
}

func (s *AppService) GetAppCount(where map[string]interface{}) (count int) {
	err := models.Count(&models.App{}, where, &count)
	if err != nil {
		logger.Error("GetAppCount Error:", err)
		return 0
	}
	return
}

func (s *AppService) GetUsageAppCount(where map[string]interface{}) (count int) {
	condition := "status != -1"
	for key, val := range where {
		condition += fmt.Sprintf(" and %s=%v", key, val)
	}
	err := models.CountByWhereString(&models.App{}, condition, &count)
	if err != nil {
		logger.Error("GetUsageAppCount Error:", err)
		return 0
	}
	return
}

func (s *AppService) GetAppPageList(where map[string]interface{}, page int, pageSize int) (list []models.App, count int) {
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
	err := models.PageListbyWhereString(&models.App{}, condition, page, pageSize, "created_at desc", &list, &count)
	if err != nil {
		logger.Error("GetAppPageList Error:", err)
		return nil, 0
	}
	return
}

func (s *AppService) GetAppByID(id int64) *models.App {
	var app models.App
	err := models.Get(id, &app)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			logger.Error("GetAppByID Error:", err)
		}
		return nil
	}
	return &app
}

func (s *AppService) GetAppByAgentID(id int64) (list []models.App) {
	err := models.Where(&list, "agent_id = ? and status != ?", id, constants.APP_STATUS_PAUSE)
	if err != nil {
		logger.Error("GetAppByAgentID Error:", err)
		return nil
	}
	return
}

func (s *AppService) GetUsageAppByAgentID(id int64) (list []models.App) {
	err := models.Where(&list,
		"agent_id = ? and status in (?,?,?)",
		id,
		constants.APP_STATUS_UNKNOWN,
		constants.APP_STATUS_RUNNING,
		constants.APP_STATUS_STOP)
	if err != nil {
		logger.Error("GetUsageAppByAgentID Error:", err)
		return nil
	}
	return
}

func (s *AppService) CreateApp(app *models.App) bool {
	err := models.Create(app)
	if err != nil {
		logger.Error("CreateApp Error:", err)
		return false
	}
	return true
}

func (s *AppService) UpdateApp(app *models.App) bool {
	err := models.Update(app)
	if err != nil {
		logger.Error("UpdateApp Error:", err)
		return false
	}
	return true
}

func (s *AppService) DeleteAppByID(id int64) bool {
	app := new(models.App)
	app.ID = id
	err := models.Delete(app)
	if err != nil {
		logger.Error("DeleteAppByID Error:", err)
		return false
	}
	return true
}

func (s *AppService) ChangeAPPStatus(app *models.App, status int64, updator int64) bool {
	values := map[string]interface{}{
		"status":  status,
		"updator": updator,
	}
	err := models.UpdateColumns(app, values)
	if err != nil {
		logger.Error("ChangeAPPStatus Error:", err)
		return false
	}
	return true
}
