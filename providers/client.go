package providers

import (
	"Asgard/clients"
	"Asgard/models"
)

var (
	AgentClients = map[int64]*clients.Agent{}
)

func GetAgent(agent *models.Agent) (*clients.Agent, error) {
	client, ok := AgentClients[agent.ID]
	if ok {
		return client, nil
	}
	client, err := clients.NewAgent(agent.IP, agent.Port)
	if err != nil {
		return nil, err
	}
	AgentClients[agent.ID] = client
	return client, nil
}
