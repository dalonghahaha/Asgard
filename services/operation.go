package services

import (
	"Asgard/models"

	"github.com/dalonghahaha/avenger/components/logger"
)

type OperationService struct {
}

func NewOperationService() *OperationService {
	return &OperationService{}
}

func (s *OperationService) GetOperationPageList(where map[string]interface{}, page int, pageSize int) (list []models.Operation, count int) {
	err := models.PageList(&models.Operation{}, where, page, pageSize, "created_at desc", &list, &count)
	if err != nil {
		logger.Error("GetOperationLogPageList Error:", err)
		return nil, 0
	}
	return
}

func (s *OperationService) CreateOperation(operationLog *models.Operation) bool {
	err := models.Create(operationLog)
	if err != nil {
		logger.Error("CreateOperationLog Error:", err)
		return false
	}
	return true
}
