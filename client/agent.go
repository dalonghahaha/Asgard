package client

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"

	"Asgard/models"
	"Asgard/rpc"
)

var (
	Agents = map[int64]rpc.AgentClient{}
)

func GetAgent(agent *models.Agent) (rpc.AgentClient, error) {
	_client, ok := Agents[agent.ID]
	if ok {
		return _client, nil
	}
	addr := fmt.Sprintf("%s:%s", agent.IP, agent.Port)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	Agents[agent.ID] = rpc.NewAgentClient(conn)
	return Agents[agent.ID], nil
}

func GetAgentStat(agent *models.Agent) (*rpc.AgentStat, error) {
	agentClient, err := GetAgent(agent)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), RPCTimeOut)
	defer cancel()
	response, err := agentClient.Stat(ctx, &rpc.Empty{})
	if err != nil {
		return nil, err
	}
	return response.GetAgentStat(), nil
}

func GetAgentAppList(agent *models.Agent) ([]*rpc.App, error) {
	agentClient, err := GetAgent(agent)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), RPCTimeOut)
	defer cancel()
	response, err := agentClient.AppList(ctx, &rpc.Empty{})
	if err != nil {
		return nil, err
	}
	return response.GetApps(), nil
}

func GetAgentApp(agent *models.Agent, id int64) (*rpc.App, error) {
	agentClient, err := GetAgent(agent)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), RPCTimeOut)
	defer cancel()
	response, err := agentClient.AppGet(ctx, &rpc.ID{Id: id})
	if err != nil {
		return nil, err
	}
	if response.GetCode() == rpc.Nofound {
		return nil, nil
	}
	return response.GetApp(), nil
}

func AddAgentApp(agent *models.Agent, app *models.App) error {
	agentClient, err := GetAgent(agent)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), RPCTimeOut)
	defer cancel()
	response, err := agentClient.AppAdd(ctx, rpc.FormatApp(app))
	if err != nil {
		return err
	}
	if response.GetCode() == rpc.OK {
		return nil
	}
	return fmt.Errorf(response.GetMessage())
}

func UpdateAgentApp(agent *models.Agent, app *models.App) error {
	agentClient, err := GetAgent(agent)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), RPCTimeOut)
	defer cancel()
	response, err := agentClient.AppUpdate(ctx, rpc.FormatApp(app))
	if err != nil {
		return err
	}
	if response.GetCode() == rpc.OK {
		return nil
	}
	return fmt.Errorf(response.GetMessage())
}

func RemoveAgentApp(agent *models.Agent, id int64) error {
	agentClient, err := GetAgent(agent)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), RPCTimeOut)
	defer cancel()
	response, err := agentClient.AppRemove(ctx, &rpc.ID{Id: id})
	if err != nil {
		return err
	}
	if response.GetCode() == rpc.OK {
		return nil
	}
	return fmt.Errorf(response.GetMessage())
}

func GetAgentJobList(agent *models.Agent) ([]*rpc.Job, error) {
	agentClient, err := GetAgent(agent)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), RPCTimeOut)
	defer cancel()
	response, err := agentClient.JobList(ctx, &rpc.Empty{})
	if err != nil {
		return nil, err
	}
	return response.GetJobs(), nil
}

func GetAgentJob(agent *models.Agent, id int64) (*rpc.Job, error) {
	agentClient, err := GetAgent(agent)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), RPCTimeOut)
	defer cancel()
	response, err := agentClient.JobGet(ctx, &rpc.ID{Id: id})
	if err != nil {
		return nil, err
	}
	if response.GetCode() == rpc.Nofound {
		return nil, nil
	}
	return response.GetJob(), nil
}

func AddAgentJob(agent *models.Agent, job *models.Job) error {
	agentClient, err := GetAgent(agent)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), RPCTimeOut)
	defer cancel()
	response, err := agentClient.JobAdd(ctx, rpc.FormatJob(job))
	if err != nil {
		return err
	}
	if response.GetCode() == rpc.OK {
		return nil
	}
	return fmt.Errorf(response.GetMessage())
}

func UpdateAgentJob(agent *models.Agent, job *models.Job) error {
	agentClient, err := GetAgent(agent)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), RPCTimeOut)
	defer cancel()
	response, err := agentClient.JobUpdate(ctx, rpc.FormatJob(job))
	if err != nil {
		return err
	}
	if response.GetCode() == rpc.OK {
		return nil
	}
	return fmt.Errorf(response.GetMessage())
}

func RemoveAgentJob(agent *models.Agent, id int64) error {
	agentClient, err := GetAgent(agent)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), RPCTimeOut)
	defer cancel()
	response, err := agentClient.JobRemove(ctx, &rpc.ID{Id: id})
	if err != nil {
		return err
	}
	if response.GetCode() == rpc.OK {
		return nil
	}
	return fmt.Errorf(response.GetMessage())
}

func GetAgentTimingList(agent *models.Agent) ([]*rpc.Timing, error) {
	agentClient, err := GetAgent(agent)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), RPCTimeOut)
	defer cancel()
	response, err := agentClient.TimingList(ctx, &rpc.Empty{})
	if err != nil {
		return nil, err
	}
	return response.GetTimings(), nil
}

func GetAgentTiming(agent *models.Agent, id int64) (*rpc.Timing, error) {
	agentClient, err := GetAgent(agent)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), RPCTimeOut)
	defer cancel()
	response, err := agentClient.TimingGet(ctx, &rpc.ID{Id: id})
	if err != nil {
		return nil, err
	}
	if response.GetCode() == rpc.Nofound {
		return nil, nil
	}
	return response.GetTiming(), nil
}

func AddAgentTiming(agent *models.Agent, timing *models.Timing) error {
	agentClient, err := GetAgent(agent)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), RPCTimeOut)
	defer cancel()
	response, err := agentClient.TimingAdd(ctx, rpc.FormatTiming(timing))
	if err != nil {
		return err
	}
	if response.GetCode() == rpc.OK {
		return nil
	}
	return fmt.Errorf(response.GetMessage())
}

func UpdateAgentTiming(agent *models.Agent, timing *models.Timing) error {
	agentClient, err := GetAgent(agent)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), RPCTimeOut)
	defer cancel()
	response, err := agentClient.TimingUpdate(ctx, rpc.FormatTiming(timing))
	if err != nil {
		return err
	}
	if response.GetCode() == rpc.OK {
		return nil
	}
	return fmt.Errorf(response.GetMessage())
}

func RemoveAgentTiming(agent *models.Agent, id int64) error {
	agentClient, err := GetAgent(agent)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), RPCTimeOut)
	defer cancel()
	response, err := agentClient.TimingRemove(ctx, &rpc.ID{Id: id})
	if err != nil {
		return err
	}
	if response.GetCode() == rpc.OK {
		return nil
	}
	return fmt.Errorf(response.GetMessage())
}
