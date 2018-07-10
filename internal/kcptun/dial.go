package kcptun

import (
	"net"
	"github.com/xtaci/kcp-go"
	"github.com/gongt/proxy-gateway/internal/constants"
	"log"
)

func DialKCP(address string, password string) net.Conn {
	encrypt := HashKcpPass(password)

	kcpconn, err := kcp.DialWithOptions(address, encrypt, constants.KcpDataShards, constants.KcpParityShards)

	if err != nil {
		log.Fatal("failed to Dial kcp: ", err)
	}

	return kcpconn
}
