package rpc_server

import (
	"google.golang.org/grpc/peer"
	"net"
	"errors"
	"context"
)

func (s *ConnectionBridgeServer) getRemoteIp(ctx context.Context) (string, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		err := errors.New("can not detect client peer info")
		return "", err
	}

	addr, err := net.ResolveTCPAddr(p.Addr.Network(), p.Addr.String())
	if err != nil {
		return "", err
	}
	return addr.IP.String(), nil
}
