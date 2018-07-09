package net_multiplex

import (
	"time"
	"context"
	"net"
)

func backgroundWithTimeout(t time.Duration) (ret context.Context) {
	return ret
}

func Dial(connect *NaiveAddr, t time.Duration) (ret net.Conn, err error) {
	ctx, cancel := context.WithCancel(context.Background())

	tmr := time.NewTimer(t)

	go func() {
		<-tmr.C
		cancel()
	}()

	var d net.Dialer
	ret, err = d.DialContext(ctx, connect.Network, connect.Address)

	tmr.Stop()

	return
}
