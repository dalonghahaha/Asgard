package clients

import (
	"context"

	"google.golang.org/grpc"

	"Asgard/constants"
	"Asgard/rpc"
)

type Guard struct {
	client rpc.GuardClient
}

func NewGuard(serverFile string) (*Guard, error) {
	conn, err := grpc.Dial(
		serverFile,
		grpc.WithInsecure(),
		grpc.WithContextDialer(UnixConnectCtx),
	)
	if err != nil {
		return nil, err
	}
	client := rpc.NewGuardClient(conn)
	agent := Guard{
		client: client,
	}
	return &agent, nil
}

func (a *Guard) GetList() ([]*rpc.App, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.RPC_TIMEOUT)
	defer cancel()
	response, err := a.client.List(ctx, &rpc.Empty{})
	if err != nil {
		return nil, err
	}
	return response.GetApps(), nil
}

func (a *Guard) Get(name string) (*rpc.App, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.RPC_TIMEOUT)
	defer cancel()
	response, err := a.client.Get(ctx, &rpc.Name{Name: name})
	if err != nil {
		return nil, err
	}
	return response.GetApp(), nil
}
