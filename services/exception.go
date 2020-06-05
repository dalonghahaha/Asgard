package services

import (
	"Asgard/models"

	"github.com/dalonghahaha/avenger/components/logger"
)

type ExceptionService struct {
}

func NewExceptionService() *ExceptionService {
	return &ExceptionService{}
}

func (s *ExceptionService) GetExceptionPageList(where map[string]interface{}, page int, pageSize int) (list []models.Exception, count int) {
	err := models.PageList(&models.Exception{}, where, page, pageSize, "created_at desc", &list, &count)
	if err != nil {
		logger.Error("GetExceptionPageList Error:", err)
		return nil, 0
	}
	return
}

func (s *ExceptionService) CreateException(exception *models.Exception) bool {
	err := models.Create(exception)
	if err != nil {
		logger.Error("CreateException Error:", err)
		return false
	}
	return true
}
