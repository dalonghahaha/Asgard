package services

import (
	"github.com/dalonghahaha/avenger/components/logger"
	"github.com/jinzhu/gorm"

	"Asgard/models"
)

type JobService struct {
}

func NewJobService() *JobService {
	return &JobService{}
}

func (s *JobService) GetJobPageList(where map[string]interface{}, page int, pageSize int) (list []models.Job, count int) {
	err := models.PageList(&models.Job{}, where, page, pageSize, &list, &count)
	if err != nil {
		logger.Error("GetJobPageList Error:", err)
		return nil, 0
	}
	return
}

func (s *JobService) GetJobByID(id int64) *models.Job {
	var job models.Job
	err := models.Get(id, &job)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			logger.Error("GetAppByID Error:", err)
		}
		return nil
	}
	return &job
}

func (s *JobService) CreateJob(job *models.Job) bool {
	err := models.Create(job)
	if err != nil {
		logger.Error("CreateApp Error:", err)
		return false
	}
	return true
}

func (s *JobService) UpdateJob(job *models.Job) bool {
	err := models.Update(job)
	if err != nil {
		logger.Error("UpdateApp Error:", err)
		return false
	}
	return true
}

func (s *JobService) DeleteJobByID(id int64) bool {
	job := new(models.Job)
	job.ID = id
	err := models.Delete(job)
	if err != nil {
		logger.Error("DeleteAppByID Error:", err)
		return false
	}
	return true
}
