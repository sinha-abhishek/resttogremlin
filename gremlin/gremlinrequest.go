package gremlin

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
)

type Arguments struct {
	Gremlin  string                 `json:"gremlin"`
	Language string                 `json:"language"`
	Bindings map[string]interface{} `json:"bindings"`
}

type GremlinRequest struct {
	RequestID string    `json:"requestId"`
	Processor string    `json:"processor"`
	Op        string    `json:"op"`
	Args      Arguments `json:"args"`
}

func (gr *GremlinRequest) SetQuery(query string) {
	gr.Args.Gremlin = query
}

func (gr *GremlinRequest) AddBinding(key string, value interface{}) {
	gr.Args.Bindings[key] = value
}

func (gr *GremlinRequest) AddBindings(keyVal map[string]interface{}) {
	for k, v := range keyVal {
		gr.Args.Bindings[k] = v
	}
}

func (gr *GremlinRequest) ReadQueryFromFileAndAddBindings(filename string, keyVal map[string]interface{}) error {
	var buffer bytes.Buffer
	gloabalVars, err1 := ioutil.ReadFile("scripts/globals.groovy")
	if err1 != nil {
		log.Println(err1)
		return err1
	}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err)
		return err
	}
	buffer.WriteString(string(gloabalVars))
	buffer.WriteString(string(data))
	gr.Args.Gremlin = buffer.String()
	gr.AddBindings(keyVal)
	return err
}

func (gr *GremlinRequest) PackageRequest() (data []byte, err error) {
	var j []byte
	j, err = json.Marshal(gr)
	if err != nil {
		return
	}

	mimetype := []byte("application/json")
	mimetypelen := byte(len(mimetype))

	data = append(data, mimetypelen)
	data = append(data, mimetype...)
	data = append(data, j...)
	return
}
