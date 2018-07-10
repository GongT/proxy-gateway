package client_config

import (
	"net/url"
	"golang.org/x/net/proxy"
	"encoding/json"
	"bytes"
	"log"
)

var GlobalDialer proxy.Dialer = proxy.Direct

type ProxySetting struct {
	str string
	url url.URL
}

var NoProxy = ProxySetting{str: ""}

func (m ProxySetting) MarshalJSON() ([]byte, error) {
	buff := bytes.NewBufferString("")
	buff.WriteByte('"')
	buff.WriteString(m.url.String())
	buff.WriteByte('"')
	return buff.Bytes(), nil
}

func (m *ProxySetting) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	u, err := url.Parse(s)
	if err != nil {
		return err
	}
	*m = ProxySetting{
		str: s,
		url: *u,
	}
	return nil
}

func ApplyProxy(setting ProxySetting, useEnv bool) (dial proxy.Dialer, err error) {
	if len(setting.str) == 0 {
		if useEnv {
			dial, err = proxy.FromEnvironment(), nil
			if dial == proxy.Direct {
				log.Println("not using a proxy.")
			} else {
				log.Println("using proxy server from environment:", setting.url.String())
			}
		} else {
			log.Println("not using a proxy.")
		}
	} else {
		dial, err = proxy.FromURL(&setting.url, proxy.Direct)
		log.Println("using proxy server:", setting.url.String())
	}
	if err != nil {
		log.Fatal("proxy setting error: ", err)
	}

	GlobalDialer = dial

	return
}
