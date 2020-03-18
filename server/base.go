package server

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

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
	// grpc_logrus.ReplaceGrpcLogger(logrus.NewEntry(grpcLogger))
	// option := grpc_middleware.WithUnaryServerChain(
	// 	recoverInterceptor(),
	// )
	return grpc.NewServer(
		grpc.MaxRecvMsgSize(1024*1024*1024),
		grpc.MaxSendMsgSize(1024*1024*1024),
	)
}

func DefaultServer() *grpc.Server {
	return grpc.NewServer()
}

func GetLog(dir string, lines int) [][]byte {
	content := [][]byte{}
	info := Shell("tail", "-"+strconv.Itoa(lines), dir)
	for _, line := range strings.Split(info, "\n") {
		content = append(content, []byte(line))
	}
	return content
}

func Shell(executor string, args ...string) string {
	cmd := exec.Command(executor, args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(err.Error(), stderr.String())
		return ""
	}
	result := out.String()
	return strings.Trim(result, "\n")
}
