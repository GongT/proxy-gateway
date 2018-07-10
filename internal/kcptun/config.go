package kcptun

import (
	"golang.org/x/crypto/pbkdf2"
	"crypto/sha1"
	"github.com/gongt/proxy-gateway/internal/constants"
	"github.com/xtaci/kcp-go"
	"log"
)

func HashKcpPass(key string) kcp.BlockCrypt {
	password := pbkdf2.Key([]byte(key), []byte(constants.Salt), 4096, 32, sha1.New)

	encrypt, err := kcp.NewTripleDESBlockCrypt(password)
	if err != nil {
		log.Fatal("failed to create kcptun crypt: ", err)
	}

	return encrypt
}

func ConfigAccept(c *kcp.UDPSession) {
	c.SetStreamMode(true)
	c.SetWriteDelay(true)
	c.SetNoDelay(1, 20, 2, 1)
	c.SetACKNoDelay(true)
}
