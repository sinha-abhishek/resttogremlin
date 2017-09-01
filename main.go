package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"bitbucket.org/abh_sinha/gremlin_client/gremlin"
	"bitbucket.org/abh_sinha/gremlin_client/handlers"
)

var handlerMap map[string]handlers.HandlerConf
var client *gremlin.Client
var configurations *Configurations

func Init() (err error) {

	configurations, err = GetConfigurations()
	if err != nil {
		log.Println(err)
		return
	}
	handlerMap, err = handlers.GetHandlerConfigs()
	if err != nil {
		fmt.Println(err)
		return
	}
	client = gremlin.NewClient(configurations.GremlinHosts)
	//err = client.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}

func test() {
	handler := handlerMap["init"]
	bMap := make(map[string]interface{})
	bMap["uid"] = 1006
	bMap["number2"] = "10030"
	bMap["contact_name"] = "testnam2"
	query := handler.Query
	bindingMap := handler.Bindings
	bindings := make(map[string]interface{})
	for k, v := range bindingMap {
		bindings[k] = bMap[v]
	}

	gr := client.NewGremlinRequest()
	gr.SetQuery(query)
	gr.AddBindings(bindings)
	s, _ := client.SendRequest(gr)
	fmt.Println(s)
}

func handleGremlinRequest(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// for key, values := range r.PostForm {
	// 	// [...]
	// 	// fmt.Println(key)
	// 	// fmt.Println(values)
	// }

	method := r.PostFormValue("method")

	fmt.Println(method)
	if handler, ok := handlerMap[method]; ok {
		query := handler.Query
		bindingMap := handler.Bindings
		bindings := make(map[string]interface{})
		requestMap := make(map[string]interface{})
		for k := range r.PostForm {
			requestMap[k] = r.PostFormValue(k)
		}
		for k, v := range bindingMap {
			if rVal, ok := requestMap[v]; ok {
				bindings[k] = rVal
			} else {
				http.Error(w, "values not sent : "+v, 400)
				return
			}
		}
		gr := client.NewGremlinRequest()
		gr.SetQuery(query)
		gr.AddBindings(bindings)
		s, err2 := client.SendRequest(gr)
		if err2 != nil {
			http.Error(w, err2.Error(), 400)
			return
		}
		fmt.Println(s)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(s); err != nil {
			http.Error(w, "server error", 500)
			panic(err)
		}
	} else {
		http.Error(w, "method not found", 400)
		return
	}

}

func StartServer(port string) {
	http.HandleFunc("/gremlin", handleGremlinRequest)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Println(err)
		panic(1)
	}
}

func main() {
	err := Init()
	if err != nil {
		return
	}
	StartServer(configurations.ServerPort)
}
