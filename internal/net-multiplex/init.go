package net_multiplex

import (
	"github.com/hashicorp/yamux"
	"os"
	"time"
	"github.com/gongt/proxy-gateway/internal/constants"
)

var keepalive = constants.DefaultKeepAlive * time.Second

func SetKeepAlive(ka time.Duration) {
	keepalive = ka
}

func init() {
	cfg := yamux.DefaultConfig()
	cfg.EnableKeepAlive = true
	cfg.KeepAliveInterval = keepalive
	cfg.ConnectionWriteTimeout = keepalive / 2
	cfg.LogOutput = os.Stderr
}
