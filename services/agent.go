package services

import (
	"github.com/dalonghahaha/avenger/components/logger"
	"github.com/jinzhu/gorm"

	"Asgard/models"
)

type AgentService struct {
}

func NewAgentService() *AgentService {
	return &AgentService{}
}

func (s *AgentService) GetAllAgent() (list []models.Agent) {
	err := models.Where(&list, "status != ?", "-1")
	if err != nil {
		logger.Error("GetAllAgent Error:", err)
		return nil
	}
	return list
}

func (s *AgentService) GetOnlineAgent() (list []models.Agent) {
	err := models.Where(&list, "status = ?", "1")
	if err != nil {
		logger.Error("GetOnlineAgent Error:", err)
		return nil
	}
	return list
}

func (s *AgentService) GetOfflineAgent() (list []models.Agent) {
	err := models.Where(&list, "status = ?", "0")
	if err != nil {
		logger.Error("GetOfflineAgent Error:", err)
		return nil
	}
	return list
}

func (s *AgentService) GetAgentPageList(where map[string]interface{}, page int, pageSize int) (list []models.Agent, count int) {
	err := models.PageList(&models.Agent{}, where, page, pageSize, &list, &count)
	if err != nil {
		logger.Error("GetAgentPageList Error:", err)
		return nil, 0
	}
	return
}

func (s *AgentService) GetAgentByID(id int64) *models.Agent {
	var agent models.Agent
	err := models.Get(id, &agent)
	if err != nil && err != gorm.ErrRecordNotFound {
		if err != gorm.ErrRecordNotFound {
			logger.Error("GetAgentByID Error:", err)
		}
		return nil
	}
	return &agent
}

func (s *AgentService) GetAgentByIPAndPort(ip, port string) *models.Agent {
	where := map[string]interface{}{
		"ip":   ip,
		"port": port,
	}
	var agent models.Agent
	err := models.Find(where, &agent)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			logger.Error("GetAgentByIP Error:", err)
		}
		return nil
	}
	return &agent
}

func (s *AgentService) CreateAgent(agent *models.Agent) bool {
	err := models.Create(agent)
	if err != nil {
		logger.Error("CreateAgent Error:", err)
		return false
	}
	return true
}

func (s *AgentService) UpdateAgent(agent *models.Agent) bool {
	err := models.Update(agent)
	if err != nil {
		logger.Error("UpdateAgent Error:", err)
		return false
	}
	return true
}

func (s *AgentService) DeleteAgentByID(id int64) bool {
	agent := new(models.Agent)
	agent.ID = id
	err := models.Delete(agent)
	if err != nil {
		logger.Error("DeleteAgentByID Error:", err)
		return false
	}
	return true
}
