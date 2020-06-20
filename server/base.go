package server

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/dalonghahaha/avenger/components/logger"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"

	"Asgard/rpc"
)

type baseServer struct{}

func (s *baseServer) OK() (*rpc.Response, error) {
	return &rpc.Response{Code: rpc.OK, Message: "ok"}, nil
}

func (s *baseServer) Error(msg string) (*rpc.Response, error) {
	return &rpc.Response{Code: rpc.Error, Message: msg}, nil
}

func accessLoggerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		service := path.Dir(info.FullMethod)[1:]
		method := path.Base(info.FullMethod)
		begin := time.Now()
		resp, err := handler(ctx, req)
		finish := time.Now()
		nanoseconds := finish.Sub(begin).Nanoseconds()
		milliseconds := fmt.Sprintf("%d.%d", nanoseconds/1e6, nanoseconds%1e6)
		logger.Debugf("rpc call:[%s][%s][%s]", service, method, milliseconds)
		return resp, err
	}
}

func NewRPCServer() *grpc.Server {
	option := grpc_middleware.WithUnaryServerChain(
		accessLoggerInterceptor(),
	)
	return grpc.NewServer(
		grpc.MaxRecvMsgSize(1024*1024*1024),
		grpc.MaxSendMsgSize(1024*1024*1024),
		option,
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
