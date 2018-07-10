package kcptun

import (
	"net"
	"github.com/gongt/proxy-gateway/internal/constants"
	"log"
)

func ListenUDP(network string) *net.UDPConn {
	udpaddr, err := net.ResolveUDPAddr(network, constants.PortCommunicate)
	if err != nil {
		log.Fatal("failed to ResolveUDPAddr: ", err)
	}
	conn, err := net.ListenUDP(network, udpaddr)
	if err != nil {
		log.Fatal("failed to ListenUDP: ", err)
	}
	return conn
}
