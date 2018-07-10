package config_file

import (
	"encoding/json"
	"log"
	"os"
	"flag"
	"github.com/gongt/proxy-gateway/internal/my-strings"
	"github.com/gongt/proxy-gateway/internal/config/client_config"
)

type ListenConfig struct {
	Net  string `json:"type"`
	Addr string `json:"address"`
}

func (m ListenConfig) Network() string {
	return m.Net
}
func (m ListenConfig) String() string {
	return m.Addr
}

type OpenPorts struct {
	Listen  ListenConfig `json:"listen"`
	Connect ListenConfig `json:"connect"`
}

type SameConfig struct {
	Kcptun string `json:"kcptun"`
}

type ClientConfig struct {
	SameConfig
	Server    string                     `json:"server"`
	OpenPorts []OpenPorts                `json:"openPorts"`
	Proxy     client_config.ProxySetting `json:"proxy"`
}

type ServerConfig struct {
	SameConfig
}

func load(ret interface{}) { // * -> ServerConfig|ClientConfig
	if !flag.Parsed() {
		flag.Parse()
	}

	if len(configFilePath) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	f, err := os.Open(configFilePath)
	if err != nil {
		log.Fatal("Cannot open config file:", err)
	}
	defer f.Close()

	j := json.NewDecoder(f)
	err = j.Decode(ret)
	if err != nil {
		log.Fatal("Config file not ok:", err)
	}
}

func LoadServerConfig() (ret ServerConfig) {
	load(&ret)
	return
}

func LoadClientConfig() (ret ClientConfig) {
	load(&ret)

	if len(ret.Server) == 0 {
		missingField("server", "where to connect")
	}
	if len(ret.OpenPorts) == 0 {
		missingField("openPorts", "port to open, see README.")
	}

	my_strings.NormalizeServer(&ret.Server)

	return
}

func missingField(field, message string) {
	log.Printf("error parsing config file: missing field `%s`: %s", field, message)
	os.Exit(2)
}
