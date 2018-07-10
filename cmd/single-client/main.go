package main

import (
	"log"
	"github.com/gongt/proxy-gateway/internal/net-multiplex/client"
	"github.com/gongt/proxy-gateway/internal/net-multiplex"
	"net"
	"flag"
	"github.com/gongt/proxy-gateway/internal/my-strings"
	"fmt"
	"os"
	"github.com/gongt/proxy-gateway/internal/kcptun"
	"github.com/gongt/proxy-gateway/internal/config/client_config"
)

func main() {
	var remote string
	var lpListen string
	var lpConnect string
	var kcpPass string

	getString(&remote, "server", "s", "proxy server location.")
	getString(&lpListen, "listen", "l", "remote listen address, Eg: tcp://0.0.0.0:6000")
	getString(&lpConnect, "connect", "t", "local connect address, Eg: unix:///var/lib/mysql.sock")
	getString(&kcpPass, "kcp", "P", "kcpPass password, must same with server.")

	flag.Parse()

	if len(remote) <= 0 {
		usage("invalid remote addr")
	}

	listen, err := net_multiplex.ParseAddress(lpListen, "0.0.0.0")
	if err != nil {
		usage("invalid argument listen: ", err)
	}
	connect, err := net_multiplex.ParseAddress(lpConnect, "127.0.0.1")
	if err != nil {
		usage("invalid argument connect: ", err)
	}

	host, port := my_strings.NormalizeServer(&remote)
	host, _ = my_strings.GetIp(host)
	remoteFull := host + ":" + port

	log.Printf("configuration: all connection to %s on the remote host %s, will pass to %s on local network.\n", listen.String(), remoteFull, connect.String())

	client_config.ApplyProxy(client_config.NoProxy, true)

	var conn net.Conn
	if len(kcpPass) == 0 {
		conn = net_multiplex.DialTCP(remoteFull)
	} else {
		conn = kcptun.DialKCP(remoteFull, kcpPass)
	}

	c := client.NewMultiplexClient(conn)

	remoteId := c.Open(listen, connect)

	log.Printf("request success, connection type is %d. waitting client...\n", remoteId)

	c.EventLoop()
}

func usage(msg ... interface{}) {
	flag.Usage()
	fmt.Print("[!!!] ")
	fmt.Println(msg...)
	os.Exit(1)
}

func getString(p *string, long, short, desc string) {
	flag.StringVar(p, long, "", desc)
	flag.StringVar(p, short, "", desc+" (--"+long+")")
}
