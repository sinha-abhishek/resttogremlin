package main

import (
	"encoding/json"
	"log"
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
		config = new(Configurations)
		file, err = os.Open("conf/server_conf.json")
		if err == nil {
			decoder := json.NewDecoder(file)
			log.Println(decoder)
			err = decoder.Decode(config)
		}
	}
	return config, err
}
