package config_file

import "flag"

var configFilePath string

func init() {
	flag.StringVar(&configFilePath, "config", "", "(required) config file location.")
	flag.StringVar(&configFilePath, "c", "", "(required) config file location. (short hand)")
}
