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
	Network string
	Address string
}

func (addr *NaiveAddr) String() string {
	return addr.Network + "://" + addr.Address
}
func (addr *NaiveAddr) valid(defHost string) error {
	switch addr.Network {
	case "tcp":
		if !testPort.MatchString(addr.Address) {
			return errors.New("tcp address must have a port")
		}

		if strings.HasPrefix(addr.Address, ":") {
			addr.Address = defHost + addr.Address
		} else if strings.HasPrefix(addr.Address, "*:") {
			addr.Address = defHost + addr.Address[1:]
		}
	case "unix":
		if !strings.HasPrefix(addr.Address, "@/") && ! strings.HasPrefix(addr.Address, "/") {
			return errors.New("unix socket address must be absolute")
		}
	default:
		return errors.New("unknown protocol: " + addr.Network)
	}

	return nil
}

func ParseAddress(str string, defHost string) (ret *NaiveAddr, err error) {
	listen := strings.Split(str, "://")

	if len(listen) == 1 {
		addr := listen[0]
		if strings.HasPrefix(addr, "@/") || strings.HasPrefix(addr, "/") {
			ret = &NaiveAddr{Network: "unix", Address: addr}
		} else {
			ret = &NaiveAddr{Network: "tcp", Address: addr}
		}
		err = ret.valid(defHost)
		return
	} else if len(listen) == 2 {
		ret = &NaiveAddr{Network: listen[0], Address: listen[1]}
		err = ret.valid(defHost)
	} else {
		err = errors.New("invalid address: " + str)
	}
	return
}
