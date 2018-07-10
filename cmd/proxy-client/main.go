package main

import (
	"github.com/gongt/proxy-gateway/internal/config/config_file"
	"github.com/gongt/proxy-gateway/internal/net-multiplex"
	"github.com/gongt/proxy-gateway/internal/net-multiplex/client"
	"log"
	"github.com/gongt/proxy-gateway/internal/systemd"
	"github.com/gongt/proxy-gateway/internal/config/client_config"
)

func main() {
	configs := config_file.LoadClientConfig()

	client_config.ApplyProxy(configs.Proxy, true)

	conn := net_multiplex.DialTCP(configs.Server)
	log.Println("connected.")
	c := client.NewMultiplexClient(conn)

	systemd.Status("connecting")

	for _, value := range configs.OpenPorts {
		remoteId := c.Open(value.Listen, value.Connect)
		log.Printf("request success, connection type is %d. waitting client...\n", remoteId)
	}
	log.Printf("All port request success!")

	systemd.Ready("listening")

	c.EventLoop()
}
