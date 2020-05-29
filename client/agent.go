package client

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"Asgard/constants"
	"Asgard/models"
	"Asgard/rpc"
)

type Agent struct {
	client rpc.AgentClient
}

func NewAgent(ip, port string) (*Agent, error) {
	addr := fmt.Sprintf("%s:%s", ip, port)
	ctx, cancel := context.WithTimeout(context.Background(), constants.RPC_TIMEOUT)
	defer cancel()
	option := grpc.WithDefaultCallOptions(
		grpc.MaxCallRecvMsgSize(constants.RPC_MESSAGE_SIZE),
		grpc.MaxCallSendMsgSize(constants.RPC_MESSAGE_SIZE),
	)
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), option)
	if err != nil {
		return nil, err
	}
	client := rpc.NewAgentClient(conn)
	agent := Agent{
		client: client,
	}
	return &agent, nil
}

func (a *Agent) GetStat() (*rpc.AgentStat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.RPC_TIMEOUT)
	defer cancel()
	response, err := a.client.Stat(ctx, &rpc.Empty{})
	if err != nil {
		return nil, err
	}
	return response.GetAgentStat(), nil
}

func (a *Agent) GetLog(agent *models.Agent, dir string, lines int64) ([]string, error) {
	content := []string{}
	ctx, cancel := context.WithTimeout(context.Background(), constants.RPC_TIMEOUT)
	defer cancel()
	response, err := a.client.Log(ctx, &rpc.LogRuquest{Dir: dir, Lines: lines})
	if err != nil {
		return content, err
	}
	for _, val := range response.GetContent() {
		content = append(content, string(val))
	}
	return content, nil
}

func (a *Agent) GetAppList(agent *models.Agent) ([]*rpc.App, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.RPC_TIMEOUT)
	defer cancel()
	response, err := a.client.AppList(ctx, &rpc.Empty{})
	if err != nil {
		return nil, err
	}
	return response.GetApps(), nil
}

func (a *Agent) GetApp(agent *models.Agent, id int64) (*rpc.App, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.RPC_TIMEOUT)
	defer cancel()
	response, err := a.client.AppGet(ctx, &rpc.ID{Id: id})
	if err != nil {
		return nil, err
	}
	if response.GetCode() == rpc.Nofound {
		return nil, nil
	}
	return response.GetApp(), nil
}

func (a *Agent) AddApp(agent *models.Agent, app *models.App) error {
	ctx, cancel := context.WithTimeout(context.Background(), constants.RPC_TIMEOUT)
	defer cancel()
	response, err := a.client.AppAdd(ctx, rpc.FormatApp(app))
	if err != nil {
		return err
	}
	if response.GetCode() == rpc.OK {
		return nil
	}
	return fmt.Errorf(response.GetMessage())
}

func (a *Agent) UpdateApp(agent *models.Agent, app *models.App) error {
	ctx, cancel := context.WithTimeout(context.Background(), constants.RPC_TIMEOUT)
	defer cancel()
	response, err := a.client.AppUpdate(ctx, rpc.FormatApp(app))
	if err != nil {
		return err
	}
	if response.GetCode() == rpc.OK {
		return nil
	}
	return fmt.Errorf(response.GetMessage())
}

func (a *Agent) RemoveApp(agent *models.Agent, id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), constants.RPC_TIMEOUT)
	defer cancel()
	response, err := a.client.AppRemove(ctx, &rpc.ID{Id: id})
	if err != nil {
		return err
	}
	if response.GetCode() == rpc.OK {
		return nil
	}
	return fmt.Errorf(response.GetMessage())
}

func (a *Agent) GetJobList(agent *models.Agent) ([]*rpc.Job, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.RPC_TIMEOUT)
	defer cancel()
	response, err := a.client.JobList(ctx, &rpc.Empty{})
	if err != nil {
		return nil, err
	}
	return response.GetJobs(), nil
}

func (a *Agent) GetJob(agent *models.Agent, id int64) (*rpc.Job, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.RPC_TIMEOUT)
	defer cancel()
	response, err := a.client.JobGet(ctx, &rpc.ID{Id: id})
	if err != nil {
		return nil, err
	}
	if response.GetCode() == rpc.Nofound {
		return nil, nil
	}
	return response.GetJob(), nil
}

func (a *Agent) AddJob(agent *models.Agent, job *models.Job) error {
	ctx, cancel := context.WithTimeout(context.Background(), constants.RPC_TIMEOUT)
	defer cancel()
	response, err := a.client.JobAdd(ctx, rpc.FormatJob(job))
	if err != nil {
		return err
	}
	if response.GetCode() == rpc.OK {
		return nil
	}
	return fmt.Errorf(response.GetMessage())
}

func (a *Agent) UpdateJob(agent *models.Agent, job *models.Job) error {
	ctx, cancel := context.WithTimeout(context.Background(), constants.RPC_TIMEOUT)
	defer cancel()
	response, err := a.client.JobUpdate(ctx, rpc.FormatJob(job))
	if err != nil {
		return err
	}
	if response.GetCode() == rpc.OK {
		return nil
	}
	return fmt.Errorf(response.GetMessage())
}

func (a *Agent) RemoveJob(agent *models.Agent, id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), constants.RPC_TIMEOUT)
	defer cancel()
	response, err := a.client.JobRemove(ctx, &rpc.ID{Id: id})
	if err != nil {
		return err
	}
	if response.GetCode() == rpc.OK {
		return nil
	}
	return fmt.Errorf(response.GetMessage())
}

func (a *Agent) GetTimingList(agent *models.Agent) ([]*rpc.Timing, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.RPC_TIMEOUT)
	defer cancel()
	response, err := a.client.TimingList(ctx, &rpc.Empty{})
	if err != nil {
		return nil, err
	}
	return response.GetTimings(), nil
}

func (a *Agent) GetTiming(agent *models.Agent, id int64) (*rpc.Timing, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.RPC_TIMEOUT)
	defer cancel()
	response, err := a.client.TimingGet(ctx, &rpc.ID{Id: id})
	if err != nil {
		return nil, err
	}
	if response.GetCode() == rpc.Nofound {
		return nil, nil
	}
	return response.GetTiming(), nil
}

func (a *Agent) AddTiming(agent *models.Agent, timing *models.Timing) error {
	ctx, cancel := context.WithTimeout(context.Background(), constants.RPC_TIMEOUT)
	defer cancel()
	response, err := a.client.TimingAdd(ctx, rpc.FormatTiming(timing))
	if err != nil {
		return err
	}
	if response.GetCode() == rpc.OK {
		return nil
	}
	return fmt.Errorf(response.GetMessage())
}

func (a *Agent) UpdateTiming(agent *models.Agent, timing *models.Timing) error {
	ctx, cancel := context.WithTimeout(context.Background(), constants.RPC_TIMEOUT)
	defer cancel()
	response, err := a.client.TimingUpdate(ctx, rpc.FormatTiming(timing))
	if err != nil {
		return err
	}
	if response.GetCode() == rpc.OK {
		return nil
	}
	return fmt.Errorf(response.GetMessage())
}

func (a *Agent) RemoveTiming(agent *models.Agent, id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), constants.RPC_TIMEOUT)
	defer cancel()
	response, err := a.client.TimingRemove(ctx, &rpc.ID{Id: id})
	if err != nil {
		return err
	}
	if response.GetCode() == rpc.OK {
		return nil
	}
	return fmt.Errorf(response.GetMessage())
}
