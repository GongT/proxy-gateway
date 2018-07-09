package net_multiplex

import (
	"github.com/hashicorp/yamux"
	"os"
	"time"
)

func init() {
	cfg := yamux.DefaultConfig()
	cfg.EnableKeepAlive = true
	cfg.KeepAliveInterval = 15 * time.Second
	cfg.ConnectionWriteTimeout = 10 * time.Second
	cfg.LogOutput = os.Stderr
}
