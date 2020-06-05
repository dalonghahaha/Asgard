package clients

import (
	"context"

	"google.golang.org/grpc"

	"Asgard/constants"
	"Asgard/rpc"
)

type Cron struct {
	client rpc.CronClient
}

func NewCron(serverFile string) (*Cron, error) {
	conn, err := grpc.Dial(
		serverFile,
		grpc.WithInsecure(),
		grpc.WithContextDialer(UnixConnectCtx),
	)
	if err != nil {
		return nil, err
	}
	client := rpc.NewCronClient(conn)
	agent := Cron{
		client: client,
	}
	return &agent, nil
}

func (a *Cron) GetList() ([]*rpc.Job, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.RPC_TIMEOUT)
	defer cancel()
	response, err := a.client.List(ctx, &rpc.Empty{})
	if err != nil {
		return nil, err
	}
	return response.GetJobs(), nil
}

func (a *Cron) Get(name string) (*rpc.Job, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.RPC_TIMEOUT)
	defer cancel()
	response, err := a.client.Get(ctx, &rpc.Name{Name: name})
	if err != nil {
		return nil, err
	}
	return response.GetJob(), nil
}
