package net_multiplex

import (
	"time"
	"net"
	"log"
	"github.com/gongt/proxy-gateway/internal/config/client_config"
	"golang.org/x/net/proxy"
	"errors"
)

func DialTCP(server string) (ret net.Conn) {
	ret, err := dial(client_config.GlobalDialer, &NaiveAddr{
		network: "tcp",
		address: server,
	}, 10*time.Second)
	if err != nil {
		log.Fatal("failed to connect to "+server+", ", err)
	}
	return
}
func Dial(server net.Addr) (ret net.Conn, err error) {
	return dial(client_config.GlobalDialer, server, 10*time.Second)
}

func dial(dialer proxy.Dialer, connect net.Addr, t time.Duration) (ret net.Conn, err error) {
	hasDone := false
	done := make(chan byte)
	tmr := time.NewTimer(t)

	log.Printf("[TCP] connecting to: %s...\n", connect.String())

	go func() {
		<-tmr.C
		if !hasDone {
			hasDone = true
			err = errors.New("timeout")
		}
		close(done)
	}()

	go func() {
		ret, err = dialer.Dial(connect.Network(), connect.String())
		hasDone = true
		close(done)
		tmr.Stop()
	}()

	<-done

	return
}
