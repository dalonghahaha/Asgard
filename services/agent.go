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

func (s *AgentService) GetAllAgent() []*models.Agent {
	list, err := new(models.Agent).All()
	if err != nil {
		logger.Error("GetAllAgent Error:", err)
		return nil
	}
	return list
}

func (s *AgentService) GetAgentByID(id int64) *models.Agent {
	agent := new(models.Agent)
	err := agent.Get(id)
	if err != nil && err != gorm.ErrRecordNotFound {
		if err != gorm.ErrRecordNotFound {
			logger.Error("GetAgentByID Error:", err)
		}
		return nil
	}
	return agent
}

func (s *AgentService) GetAgentByIP(ip string) *models.Agent {
	where := map[string]interface{}{
		"ip": ip,
	}
	agent := new(models.Agent)
	err := agent.Find(where)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			logger.Error("GetAgentByIP Error:", err)
		}
		return nil
	}
	return agent
}

func (s *AgentService) CreateAgent(agent *models.Agent) bool {
	err := agent.Create()
	if err != nil {
		logger.Error("CreateAgent Error:", err)
		return false
	}
	return true
}

func (s *AgentService) UpdateAgent(agent *models.Agent) bool {
	err := agent.Update()
	if err != nil {
		logger.Error("UpdateAgent Error:", err)
		return false
	}
	return true
}

func (s *AgentService) DeleteAgentByID(id int64) bool {
	agent := new(models.Agent)
	agent.ID = id
	err := agent.Delete()
	if err != nil {
		logger.Error("DeleteAgentByID Error:", err)
		return false
	}
	return true
}
