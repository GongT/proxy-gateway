package main

import (
	"os"
	"github.com/gongt/proxy-gateway/internal/constants"
	"log"
	"github.com/gongt/proxy-gateway/internal/net-multiplex/client"
	"strconv"
	"github.com/gongt/proxy-gateway/internal/net-multiplex"
	"time"
	"net"
)

func main() {
	log.Println(os.Args)
	if len(os.Args) != 4 { // with $0
		usage("invalid argument count (" + strconv.Itoa(len(os.Args)) + ")")
	}

	remote := os.Args[1]
	if len(remote) <= 0 {
		usage("invalid remote addr")
	}

	addr := net.ParseIP(remote)
	if addr == nil {
		ips, err := net.LookupIP(remote)
		if err != nil || len(ips) == 0 {
			log.Fatal("Cannot resolve host: " + remote)
		}

		log.Printf("resolved host %s IP: %s", remote, ips)
		log.Println("using first one.")
		addr = ips[0]
	}

	listen, err := net_multiplex.ParseAddress(os.Args[2], "0.0.0.0")
	if err != nil {
		log.Fatal("invalid argument 2: ", err)
	}
	connect, err := net_multiplex.ParseAddress(os.Args[3], "127.0.0.1")
	if err != nil {
		log.Fatal("invalid argument 3: ", err)
	}

	remoteFull := addr.String() + constants.PortCommunicate

	log.Printf("configuration: all connection to %s on the remote host %s, will pass to %s on local network.\n", listen.String(), remoteFull, connect.String())

	log.Printf("connecting to: %s...\n", remoteFull)

	conn, err := net_multiplex.Dial(&net_multiplex.NaiveAddr{
		Network: "tcp",
		Address: remoteFull,
	}, 10*time.Second)
	if err != nil {
		log.Fatal("failed to connect to "+remoteFull+", ", err)
	}

	c := client.NewMultiplexClient(conn)

	var remoteId uint32
	if listen.Network == "tcp" {
		log.Println("connecting with OpenTCP.")
		remoteId = c.OpenTCP(listen.Address, connect)
	} else {
		log.Println("connecting with OpenUnix.")
		remoteId = c.OpenUnix(listen.Address, connect)
	}

	log.Printf("request success, connection type is %d. waitting client...\n", remoteId)

	c.EventLoop()
}

func usage(err string) {
	log.Fatal(err, "\n\tUsage: $0 <server ip> RemoteListen([unix://]<listen>) LocalConnect([unix://]<listen>)")
}
