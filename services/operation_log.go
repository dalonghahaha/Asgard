package services

import (
	"Asgard/models"

	"github.com/dalonghahaha/avenger/components/logger"
)

type OperationLogService struct {
}

func NewOperationLogService() *OperationLogService {
	return &OperationLogService{}
}

func (s *OperationLogService) GetOperationLogPageList(where map[string]interface{}, page int, pageSize int) (list []models.OperationLog, count int) {
	err := models.PageList(&models.OperationLog{}, where, page, pageSize, "created_at desc", &list, &count)
	if err != nil {
		logger.Error("GetOperationLogPageList Error:", err)
		return nil, 0
	}
	return
}

func (s *OperationLogService) CreateOperationLog(operationLog *models.OperationLog) bool {
	err := models.Create(operationLog)
	if err != nil {
		logger.Error("CreateOperationLog Error:", err)
		return false
	}
	return true
}
