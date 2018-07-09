package main

import (
	"net"
	"github.com/gongt/proxy-gateway/internal/constants"
	"log"
	"github.com/gongt/proxy-gateway/internal/net-multiplex/server"
)

var ch chan net.Conn

func main() {
	ch = make(chan net.Conn)
	log.Println("server starting ...")
	go doListen("tcp4")
	go doListen("tcp6")

	for c := range ch {
		handle(c)
	}
}

func doListen(network string) {
	ln, err := net.Listen(network, constants.PortCommunicate)

	if err != nil {
		log.Fatal("failed to listen on main port: ", err)
	}

	log.Println("listening on", ln.Addr().String(), "...")

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
