package client

import (
	"context"
	"fmt"

	"github.com/dalonghahaha/avenger/components/logger"
	"google.golang.org/grpc"

	"Asgard/models"
	"Asgard/rpc"
)

var (
	GuardClients = map[int64]rpc.GuardClient{}
)

func GetGuardClient(agent *models.Agent) (rpc.GuardClient, error) {
	_client, ok := GuardClients[agent.ID]
	if ok {
		return _client, nil
	}
	addr := fmt.Sprintf("%s:%s", agent.IP, agent.Port)
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		logger.Error(fmt.Sprintf("%s is offline:%v", addr, err))
		return nil, err
	}
	guardClient := rpc.NewGuardClient(conn)
	GuardClients[agent.ID] = guardClient
	return guardClient, nil
}

func GetGuardList(agent *models.Agent) ([]*rpc.App, error) {
	guardClient, err := GetGuardClient(agent)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	defer cancel()
	response, err := guardClient.List(ctx, &rpc.Empty{})
	if err != nil {
		logger.Error(fmt.Sprintf("get guard list failed：%s", err.Error()))
		return nil, err
	}
	return response.GetApps(), nil
}
