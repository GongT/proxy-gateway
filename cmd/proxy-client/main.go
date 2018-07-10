package main

import (
	"github.com/gongt/proxy-gateway/internal/config/config_file"
	"github.com/gongt/proxy-gateway/internal/net-multiplex"
	"github.com/gongt/proxy-gateway/internal/net-multiplex/client"
	"log"
	"github.com/gongt/proxy-gateway/internal/systemd"
	"github.com/gongt/proxy-gateway/internal/config/client_config"
	"net"
	"github.com/gongt/proxy-gateway/internal/kcptun"
	"github.com/gongt/proxy-gateway/internal/my-strings"
)

func main() {
	configs := config_file.LoadClientConfig()

	client_config.ApplyProxy(configs.Proxy, true)

	host, port := my_strings.NormalizeServer(&configs.Server)
	host, _ = my_strings.GetIp(host)
	serverAddr := host + ":" + port

	var conn net.Conn
	if len(configs.Kcptun) == 0 {
		conn = net_multiplex.DialTCP(serverAddr)
	} else {
		conn = kcptun.DialKCP(serverAddr, configs.Kcptun)
	}
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
