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
	url *url.URL
	env bool
}

var NoProxy = ProxySetting{str: ""}

func (m ProxySetting) MarshalJSON() ([]byte, error) {
	if m.url != nil {
		buff := bytes.NewBufferString("")
		buff.WriteByte('"')
		buff.WriteString(m.str)
		buff.WriteByte('"')
		return buff.Bytes(), nil
	}
	if m.env {
		return []byte("true"), nil
	} else {
		return []byte("false"), nil
	}
}

func (m *ProxySetting) UnmarshalJSON(b []byte) error {
	var s = ""
	var e = true
	var u *url.URL = nil

	err := json.Unmarshal(b, &e)
	if err == nil {
		*m = ProxySetting{
			str: s,
			url: u,
			env: e,
		}
		return nil
	} else {
		if _, ok := err.(*json.UnmarshalTypeError); ok {
			err = nil
		} else {
			return err
		}
	}

	err = json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	if len(s) > 0 {
		u, err = url.Parse(s)
		if err != nil {
			return err
		}
	}

	*m = ProxySetting{
		str: s,
		url: u,
		env: e,
	}

	return nil
}

func ApplyProxy(setting ProxySetting, useEnv bool) {
	if setting.url == nil {
		if useEnv && setting.env {
			GlobalDialer = proxy.FromEnvironment()
			if GlobalDialer == proxy.Direct {
				log.Println("not using a proxy (no env set).")
			} else {
				log.Println("using proxy server from environment:", setting.url.String())
			}
		} else {
			log.Println("not using a proxy (disabled in config).")
		}
	} else {
		var err error
		GlobalDialer, err = proxy.FromURL(setting.url, proxy.Direct)
		if err != nil {
			log.Fatal("proxy setting error: ", err)
		}
		log.Println("using proxy server:", setting.url.String())
	}
}
