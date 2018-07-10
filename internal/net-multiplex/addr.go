package net_multiplex

import (
	"strings"
	"errors"
	"regexp"
)

var testPort *regexp.Regexp

func init() {
	testPort, _ = regexp.Compile(":\\d+$")
}

type NaiveAddr struct {
	network string
	address string
}

func (addr *NaiveAddr) Network() string {
	return addr.network
}
func (addr *NaiveAddr) FullString() string {
	return addr.network + "://" + addr.address
}
func (addr *NaiveAddr) String() string {
	return addr.address
}
func (addr *NaiveAddr) valid(defHost string) error {
	switch addr.network {
	case "tcp":
		if !testPort.MatchString(addr.address) {
			return errors.New("tcp address must have a port")
		}

		if strings.HasPrefix(addr.address, ":") {
			addr.address = defHost + addr.address
		} else if strings.HasPrefix(addr.address, "*:") {
			addr.address = defHost + addr.address[1:]
		}
	case "unix":
		if !strings.HasPrefix(addr.address, "@/") && ! strings.HasPrefix(addr.address, "/") {
			return errors.New("unix socket address must be absolute")
		}
	default:
		return errors.New("unknown protocol: " + addr.network)
	}

	return nil
}

func ParseAddress(str string, defHost string) (ret *NaiveAddr, err error) {
	listen := strings.Split(str, "://")

	if len(listen) == 1 {
		addr := listen[0]
		if strings.HasPrefix(addr, "@/") || strings.HasPrefix(addr, "/") {
			ret = &NaiveAddr{network: "unix", address: addr}
		} else {
			ret = &NaiveAddr{network: "tcp", address: addr}
		}
		err = ret.valid(defHost)
		return
	} else if len(listen) == 2 {
		ret = &NaiveAddr{network: listen[0], address: listen[1]}
		err = ret.valid(defHost)
	} else {
		err = errors.New("invalid address: " + str)
	}
	return
}
