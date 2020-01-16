package server

import (
	"context"
	"fmt"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"Asgard/rpc"
)

var (
	grpcLogger *logrus.Logger
)

type baseServer struct{}

func (s *baseServer) OK() (*rpc.Response, error) {
	return &rpc.Response{Code: rpc.OK, Message: "ok"}, nil
}

func (s *baseServer) Error(msg string) (*rpc.Response, error) {
	return &rpc.Response{Code: rpc.Error, Message: msg}, nil
}

func recoverInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (res interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				res = &rpc.Response{Code: 200, Message: fmt.Sprintf("%v", r)}
			}
		}()
		return handler(ctx, req)
	}
}

func NewRPCServer() *grpc.Server {
	grpc_logrus.ReplaceGrpcLogger(logrus.NewEntry(grpcLogger))
	option := grpc_middleware.WithUnaryServerChain(
		recoverInterceptor(),
	)
	return grpc.NewServer(option)
}

func DefaultServer() *grpc.Server {
	return grpc.NewServer()
}
