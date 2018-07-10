package rpc_server

import (
	"github.com/gongt/proxy-gateway/api"
	"context"
	"net"
	"log"
)

func (s *ConnectionBridgeServer) KeepAlive(ctx context.Context, msg *bridge_api_call.Empty) (ret *bridge_api_call.Empty, err error) {
	return &bridge_api_call.Empty{}, nil
}

func (s *ConnectionBridgeServer) OpenTCP(ctx context.Context, msg *bridge_api_call.OpenMessage) (ret *bridge_api_call.ProtoId, err error) {
	return s.listen("tcp", msg.Address)
}

func (s *ConnectionBridgeServer) OpenUnix(ctx context.Context, msg *bridge_api_call.OpenMessage) (ret *bridge_api_call.ProtoId, err error) {
	return s.listen("unix", msg.Address)
}

func (s *ConnectionBridgeServer) listen(socketType string, address string) (ret *bridge_api_call.ProtoId, err error) {
	ln, err := net.Listen(socketType, address)
	if err != nil {
		log.Printf("listen fail on %s:%s, error: %s", socketType, address, err.Error())
		return
	}

	s.guid++
	guid := s.guid

	log.Printf("listening on %s:%s, id: %d", socketType, address, guid)

	s.mapper[guid] = ln

	go func() {
		defer func() {
			log.Printf("defered: stop listen on %s:%s, id: %d", socketType, address, guid)
			ln.Close()
			delete(s.mapper, guid)
		}()
		for {
			conn, err := ln.Accept()

			if err != nil {
				log.Println("failed to accept:", err)

				return
			}

			log.Printf("new connection to: %s:%s, id: %d", socketType, address, guid)
			s.request <- ConnectionTarget{
				Conn: conn,
				Guid: guid,
			}
		}
	}()

	ret = &bridge_api_call.ProtoId{
		Id: uint32(guid),
	}

	return
}
