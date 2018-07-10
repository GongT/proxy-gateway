package main

import (
	"net"
	"github.com/gongt/proxy-gateway/internal/constants"
	"log"
	"github.com/gongt/proxy-gateway/internal/net-multiplex/server"
	"github.com/gongt/proxy-gateway/internal/prof"
	"github.com/gongt/proxy-gateway/internal/systemd"
	"github.com/gongt/proxy-gateway/internal/config/config_file"
	"github.com/xtaci/kcp-go"
	"github.com/gongt/proxy-gateway/internal/kcptun"
)

var ch chan net.Conn


func main() {
	configs := config_file.LoadServerConfig()

	ch = make(chan net.Conn)
	log.Println("server starting ...")

	if len(configs.Kcptun) == 0 {
		go doListen("tcp4")
		go doListen("tcp6")
	} else {
		block := kcptun.HashKcpPass(configs.Kcptun)
		go doListenKCP("udp4", block)
		go doListenKCP("udp6", block)
	}

	go prof.RunForever()
	go prof.Debug()

	systemd.Ready("listening")

	for c := range ch {
		handle(c)
	}
}

func doListenKCP(network string, encrypt kcp.BlockCrypt) {
	conn := kcptun.ListenUDP(network)

	ln, err := kcp.ServeConn(encrypt, constants.KcpDataShards, constants.KcpParityShards, conn)

	if err != nil {
		log.Fatal("failed to listen KCP: ", err)
	}

	log.Println("listening KCP on", ln.Addr().String(), "...")

	prof.Snapshot("boot")

	for {
		c, err := ln.AcceptKCP()
		kcptun.ConfigAccept(c)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("new base connection from " + c.RemoteAddr().String() + ".")
		ch <- c
	}
}
func doListen(network string) {
	ln, err := net.Listen(network, constants.PortCommunicate)
	if err != nil {
		log.Fatal("failed to listen on main port: ", err)
	}

	log.Println("listening TCP on", ln.Addr().String(), "...")

	prof.Snapshot("boot")

	for {
		c, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("new base connection from " + c.RemoteAddr().String() + ".")
		ch <- c
	}
}

func handle(c net.Conn) {
	multi, err := server.NewMultiplexer(c)
	if err != nil {
		c.Close()
		log.Println("failed handle this connection:", err)
	} else {
		log.Println("multiplexer started to handle connection")
		go multi.Handle()
	}
}
