package clients

import (
	"context"
	"net"
	"time"
)

func UnixConnect(serverFile string, t time.Duration) (net.Conn, error) {
	unixAddr, err := net.ResolveUnixAddr("unix", serverFile)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialUnix("unix", nil, unixAddr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func UnixConnectCtx(ctx context.Context, serverFile string) (net.Conn, error) {
	unixAddr, err := net.ResolveUnixAddr("unix", serverFile)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialUnix("unix", nil, unixAddr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
