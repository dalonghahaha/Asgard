package providers

import (
	"Asgard/client"
	"Asgard/constants"
	"Asgard/models"
)

var (
	MasterClient *client.Master
	AgentClients = map[int64]*client.Agent{}
)

func RegisterMaster() error {
	masterClient, err := client.NewMaster(constants.MASTER_IP, constants.MASTER_PORT)
	if err != nil {
		return err
	}
	MasterClient = masterClient
	go MasterClient.Report()
	return nil
}

func GetAgent(agent *models.Agent) (*client.Agent, error) {
	_client, ok := AgentClients[agent.ID]
	if ok {
		return _client, nil
	}
	client, err := client.NewAgent(agent.IP, agent.Port)
	if err != nil {
		return nil, err
	}
	AgentClients[agent.ID] = client
	return client, nil
}
