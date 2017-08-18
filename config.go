package main

import (
	"encoding/json"
	"os"
)

type Configurations struct {
	GremlinHosts []string `json:"gremlin_hosts"`
	ServerPort   string   `json:"rest_server_port"`
}

var config *Configurations

func GetConfigurations() (*Configurations, error) {
	var err error
	var file *os.File
	if config == nil {
		file, err = os.Open("conf/server_conf.json")
		if err == nil {
			decoder := json.NewDecoder(file)
			err = decoder.Decode(config)
		}
	}
	return config, err
}
