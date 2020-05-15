package services

import (
	"fmt"

	"github.com/dalonghahaha/avenger/components/logger"
	"github.com/jinzhu/gorm"

	"Asgard/constants"
	"Asgard/models"
)

type JobService struct {
}

func NewJobService() *JobService {
	return &JobService{}
}

func (s *JobService) GetJobCount(where map[string]interface{}) (count int) {
	err := models.Count(&models.Job{}, where, &count)
	if err != nil {
		logger.Error("GetJobCount Error:", err)
		return 0
	}
	return
}

func (s *JobService) GetJobPageList(where map[string]interface{}, page int, pageSize int) (list []models.Job, count int) {
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
	err := models.PageListbyWhereString(&models.Job{}, condition, page, pageSize, "created_at desc", &list, &count)
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
			logger.Error("GetJobByID Error:", err)
		}
		return nil
	}
	return &job
}

func (s *JobService) GetJobByAgentID(id int64) (list []models.Job) {
	err := models.Where(&list, "agent_id = ? and status != ?", id, constants.JOB_STATUS_PAUSE)
	if err != nil {
		logger.Error("GetJobByAgentID Error:", err)
		return nil
	}
	return
}

func (s *JobService) GetUsageJobByAgentID(id int64) (list []models.Job) {
	err := models.Where(
		&list,
		"agent_id = ? and status in (?,?,?)",
		id,
		constants.JOB_STATUS_UNKNOWN,
		constants.JOB_STATUS_RUNNING,
		constants.JOB_STATUS_STOP)
	if err != nil {
		logger.Error("GetUsageJobByAgentID Error:", err)
		return nil
	}
	return
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
