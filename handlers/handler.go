package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type HandlerFormat struct {
	HandlerConfs map[string]HandlerConf `json:"handlers"`
}

type HandlerConf struct {
	Query    string            `json:"query"`
	Bindings map[string]string `json:"bindings"`
}

func GetHandlerConfigs() (map[string]HandlerConf, error) {
	confData, err := ioutil.ReadFile("conf/handler_conf.json")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var data HandlerFormat
	err = json.Unmarshal(confData, &data)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return data.HandlerConfs, err
}
