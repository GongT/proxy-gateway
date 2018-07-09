package main

import (
	"net"
	"fmt"
	"io"
	"os"
)

func main() {
	c, _ := net.DialTCP("tcp", nil, &net.TCPAddr{
		IP:   net.IP{127, 0, 0, 1},
		Port: 10000,
	})

	fmt.Fprintln(c, "baaaaaaaaaaaaaaaaaaaaaaaaab")
	io.Copy(os.Stdout, c)
}
