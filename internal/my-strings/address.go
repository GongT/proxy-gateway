package my_strings

import (
	"strings"
	"github.com/gongt/proxy-gateway/internal/constants"
)

func NormalizeServer(server *string) (string, string) {
	if strings.Count(*server, ":") == 0 {
		*server += constants.PortCommunicate
	}
	ret := strings.Split(*server, ":")
	return ret[0], ret[1]
}
