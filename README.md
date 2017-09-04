# resttogremlin
This is a gremlin client written in go with a wrapper to expose gremlin requests to rest client in easiest way possible. This connects to any graph database that supports TinkerPop3 using Go and also opens a http server for client communication. TinkerPop3 uses Gremlin Server to communicate with clients using either WebSockets or REST API. Here we talk to Gremlin Server using WebSockets.
There is handler_config which helps you to define handlers and the post form parameters and corresponding gremlin script.

## Installing the server (TBD: more explanation)

go get github.com/sinha-abhishek/resttogremlin

## Defining api end points (TBD: admin console to define new APIs)
* open handler_config.json
* Let's say you want to expose an API called create_vertex with property uid with value being user_id passed by client
* add this to json array 
```
"create_vertex" : {
		"query" : "graph.addVertex(\"user\").property('uid',XUID)",
		"bindings" : {
			"XUID" : "uid"
		}
	},
 ```
 * Save this file and start the server doing `go build` and running the executable. You can set the port number to run in server_conf.json.
 * This will now provide a api way to create a vertex with user id
 * To use this api you can send 
 `` curl -X POST -vid 'method=create_vertex&uid=1' -H 'Content-Type:application/x-www-form-urlencoded' http://<server:port>/gremlin ``
 
 ## Using this as gremlin client library
 If you do not want to use it as a rest server but as a gremlin library for your go project. You can do the following
 * Include the project in your go path by either cloning or `go get`
 * Creating a Cient
 `` client := gremlin.NewClient("<array of strings with each string pointing to a gremlin server url>" ``
 For example
  ```
  hosts := [1]string {"localhost:8182"}
  client := gremlin.NewClient(hosts)
  ```
  * client maintains a connection pool of all gremlin websockets.
  * To send a request create a gremlin request you need to create `GremlinRequest` object and use `setQuery` to defing your gremlin query. You can 
  define bindings using `SetBindings`
  
  ```
  gr := client.NewGremlinRequest()
  query := "g.V().has('uid',XUID)"
  bindings := make(map[string]interface{})
  bindings['XUID'] = 1
	gr.SetQuery(query)
	gr.AddBindings(bindings)
		response, err2 := client.SendRequest(gr)
		if err2 != nil {
			//handle error
		}
 ```
 
 * There are also other methods like reading a gremlin request query from a file. Please see gremlinrequest.go
 * Response is the result returned by gremlin server

  
  
  

