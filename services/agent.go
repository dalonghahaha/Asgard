package services

import (
	"encoding/json"
	"fmt"

	"github.com/dalonghahaha/avenger/components/logger"
	"github.com/jinzhu/gorm"

	"Asgard/constants"
	"Asgard/models"
)

type AgentService struct {
}

func NewAgentService() *AgentService {
	return &AgentService{}
}

func (s *AgentService) buidCacheKey(id int64) string {
	return fmt.Sprintf("%s:%d", constants.CACHE_KEY_AGENT, id)
}

func (s *AgentService) buidCacheKeyIpPort(ip, port string) string {
	return fmt.Sprintf("%s:%s:%s", constants.CACHE_KEY_AGENT_IP_PORT, ip, port)
}

func (s *AgentService) GetUsageAgent() (list []models.Agent) {
	err := models.Where(&list, "status != ?", "-1")
	if err != nil {
		logger.Error("GetUsageAgent Error:", err)
		return nil
	}
	return list
}

func (s *AgentService) GetMasterAgent() (list []models.Agent) {
	err := models.Where(&list, "status != ? and master = ?", "-1", constants.MASTER_IP)
	if err != nil {
		logger.Error("GetUsageAgent Error:", err)
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

func (s *AgentService) GetAgentCount(where map[string]interface{}) (count int) {
	err := models.Count(&models.Agent{}, where, &count)
	if err != nil {
		logger.Error("GetAgentCount Error:", err)
		return 0
	}
	return
}

func (s *AgentService) GetAgentPageList(where map[string]interface{}, page int, pageSize int) (list []models.Agent, count int) {
	condition := "1=1"
	for key, val := range where {
		if key == "status" {
			if val.(int) == -99 {
				condition += " and status != -1"
			} else {
				condition += fmt.Sprintf(" and %s=%v", key, val)
			}
		} else if key == "alias" {
			condition += fmt.Sprintf(" and %s like '%%%v%%' ", key, val)
		} else {
			condition += fmt.Sprintf(" and %s=%v", key, val)
		}
	}
	err := models.PageListbyWhereString(&models.Agent{}, condition, page, pageSize, "name asc", &list, &count)
	if err != nil {
		logger.Error("GetAgentPageList Error:", err)
		return nil, 0
	}
	return
}

func (s *AgentService) GetAgentByID(id int64) *models.Agent {
	var agent models.Agent
	data := GetCache(s.buidCacheKey(id))
	if len(data) > 0 {
		err := json.Unmarshal([]byte(data), &agent)
		if err == nil {
			return &agent
		}
	}
	err := models.Get(id, &agent)
	if err != nil && err != gorm.ErrRecordNotFound {
		if err != gorm.ErrRecordNotFound {
			logger.Error("GetAgentByID Error:", err)
		}
		return nil
	}
	info, _ := json.Marshal(agent)
	SetCache(s.buidCacheKey(id), string(info))
	return &agent
}

func (s *AgentService) GetAgentByIPAndPort(ip, port string) *models.Agent {
	where := map[string]interface{}{
		"ip":   ip,
		"port": port,
	}
	var agent models.Agent
	data := GetCache(s.buidCacheKeyIpPort(ip, port))
	if len(data) > 0 {
		err := json.Unmarshal([]byte(data), &agent)
		if err == nil {
			return &agent
		}
	}

	err := models.Find(where, &agent)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			logger.Error("GetAgentByIP Error:", err)
		}
		return nil
	}
	info, _ := json.Marshal(agent)
	SetCache(s.buidCacheKeyIpPort(ip, port), string(info))
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
	DelCache(s.buidCacheKey(agent.ID))
	DelCache(s.buidCacheKeyIpPort(agent.IP, agent.Port))
	return true
}

func (s *AgentService) DeleteAgentByID(agent *models.Agent) bool {
	err := models.Delete(agent)
	if err != nil {
		logger.Error("DeleteAgentByID Error:", err)
		return false
	}
	DelCache(s.buidCacheKey(agent.ID))
	DelCache(s.buidCacheKeyIpPort(agent.IP, agent.Port))
	return true
}
