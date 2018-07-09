package server

import (
	"net"
	"github.com/hashicorp/yamux"
	"github.com/gongt/proxy-gateway/internal/rpc-server"
	_ "github.com/gongt/proxy-gateway/internal/net-multiplex"
	"log"
	"github.com/gongt/proxy-gateway/internal/net-multiplex"
	"encoding/binary"
	"fmt"
)

type MultiplexServer struct {
	session   *yamux.Session
	rpcServer *rpc_server.ConnectionBridgeServer
}

func NewMultiplexer(conn net.Conn) (ret *MultiplexServer, err error) {
	session, err := yamux.Server(conn, nil)
	if err != nil {
		return
	}

	ret = &MultiplexServer{
		session: session,
	}

	return
}

func (m *MultiplexServer) Handle() {
	server := rpc_server.NewRpcServer(m.session)
	m.rpcServer = server

	go func() {
		accept := server.WaitRequest()
		for conn := range accept {
			m.handleChannel(conn)
		}
	}()

	server.Start()
}

func (m *MultiplexServer) handleChannel(target rpc_server.ConnectionTarget) {
	conn, err := m.session.Open()
	if err != nil {
		log.Println("open sub connection failed:", err)
		return
	}

	log.Printf("guid: %v", uint32(target.Guid))
	err = binary.Write(conn, binary.LittleEndian, uint32(target.Guid))
	if err != nil {
		fmt.Fprintln(conn, "write id failed:", err)
		conn.Close()
		return
	}

	log.Println("open sub connection success. start bridging.")
	net_multiplex.BridgeConnectionSync(target.Conn, "accepted", conn, "remote")
}
