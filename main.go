package main

import (
	"fmt"
	"html"
	"log"
	"net/http"

	"bitbucket.org/abh_sinha/gremlin_client/gremlin"

	"github.com/gorilla/mux"
)

func main() {
	client := gremlin.NewClient("localhost:8182")
	err := client.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	gr := client.NewGremlinRequest()
	gr.SetQuery("g.V().has('uid',x).values()")
	gr.AddBinding("x", 1)
	data, err2 := gr.PackageRequest()
	fmt.Println(string(data))
	if err2 != nil {
		fmt.Println(err2)
		return
	}
	s, _ := client.SendRequest(gr)
	//s := client.SendRequest("{\"requestId\":\"655BD810-B41E-429D-B78F-3CC5F3B8E9BB\",\"processor\":\"\",\"op\":\"eval\",\"args\":{\"gremlin\":\"g.V().has('uid',1).values()\",\"language\":\"gremlin-groovy\"}}", "655BD810-B41E-429D-B78F-3CC5F3B8E9BA")
	fmt.Println(s)
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	log.Fatal(http.ListenAndServe(":8080", router))
}

/*
this is comment
*/
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}
