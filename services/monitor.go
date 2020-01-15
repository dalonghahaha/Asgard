package services

import (
	"github.com/dalonghahaha/avenger/components/logger"

	"Asgard/models"
)

type MonitorService struct {
}

func NewMonitorService() *MonitorService {
	return &MonitorService{}
}

func (s *MonitorService) GetMonitorPageList(where map[string]interface{}, page int, pageSize int) (list []models.Monitor, count int) {
	err := models.PageList(&models.Monitor{}, where, page, pageSize, &list, &count)
	if err != nil {
		logger.Error("GetMonitorPageList Error:", err)
		return nil, 0
	}
	return
}

func (s *MonitorService) CreateMonitor(monitor *models.Monitor) bool {
	err := models.Create(monitor)
	if err != nil {
		logger.Error("CreateMonitor Error:", err)
		return false
	}
	return true
}
