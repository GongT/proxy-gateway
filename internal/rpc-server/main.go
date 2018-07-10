package rpc_server

import (
	"github.com/gongt/proxy-gateway/api"
	"google.golang.org/grpc"
	"net"
	"log"
	"github.com/gongt/proxy-gateway/internal/constants"
)

type ConnectionBridgeServer struct {
	guid    uint32
	request chan ConnectionTarget

	conn   net.Listener
	rpc    *grpc.Server
	mapper map[uint32]net.Listener
	done   chan byte
}

type ConnectionTarget struct {
	Conn net.Conn
	Guid uint32
}

func NewRpcServer(conn net.Listener) (*ConnectionBridgeServer) {
	rpcServer := grpc.NewServer()

	done := make(chan byte)

	ret := ConnectionBridgeServer{
		guid:    constants.ServiceIdBase,
		request: make(chan ConnectionTarget),

		conn:   conn,
		rpc:    rpcServer,
		mapper: make(map[uint32]net.Listener),
		done:   done,
	}

	bridge_api_call.RegisterConnectionBridgeServer(rpcServer, &ret)

	return &ret
}

func (s *ConnectionBridgeServer) Start() {
	log.Println("rpc server started.")
	err := s.rpc.Serve(s.conn) // block until error
	if err != nil {
		log.Println("rpc server return with error:", err)
	} else {
		log.Println("rpc server stopped.")
	}
	s.rpc.Stop()
	s.Close()
	log.Println("rpc server finished.")
}

func (s *ConnectionBridgeServer) Close() {
	log.Printf("stop listen.")

	for _, ln := range s.mapper {
		log.Printf("stop listen on %s.", ln.Addr().String())
		ln.Close()
	}

	s.rpc.Stop()
	s.conn.Close()

	close(s.done)
	close(s.request)
}

func (s *ConnectionBridgeServer) Done() <-chan byte {
	return s.done
}

func (s *ConnectionBridgeServer) WaitRequest() <-chan ConnectionTarget {
	return s.request
}
