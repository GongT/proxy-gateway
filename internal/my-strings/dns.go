package my_strings

import (
	"net"
	"log"
	"strings"
)

func GetIp(host string) (string, string) {
	parts := strings.Split(host, ":")

	host = parts[0]

	addr := net.ParseIP(host)
	if addr == nil {
		ips, err := net.LookupIP(host)
		if err != nil || len(ips) == 0 {
			log.Fatal("Cannot resolve host: " + host)
		}

		log.Printf("resolved host %s IP: %s", host, ips)
		log.Println("using first one.")
		addr = ips[0]
		parts[0] = addr.String()
	}

	if len(parts) > 1 {
		return parts[0], parts[1]
	} else {
		return parts[0], ""
	}
}
