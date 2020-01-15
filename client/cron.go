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
	CronClients = map[int64]rpc.CronClient{}
)

func GetCronClient(agent *models.Agent) (rpc.CronClient, error) {
	_client, ok := CronClients[agent.ID]
	if ok {
		return _client, nil
	}
	addr := fmt.Sprintf("%s:%s", agent.IP, agent.Port)
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		logger.Error(fmt.Sprintf("%s is offline:%v", addr, err))
		return nil, err
	}
	cronClient := rpc.NewCronClient(conn)
	CronClients[agent.ID] = cronClient
	return cronClient, nil
}

func GetCronList(agent *models.Agent) ([]*rpc.Job, error) {
	cronClient, err := GetCronClient(agent)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	defer cancel()
	response, err := cronClient.List(ctx, &rpc.Empty{})
	if err != nil {
		logger.Error(fmt.Sprintf("get cron list failed：%s", err.Error()))
		return nil, err
	}
	return response.GetJobs(), nil
}
