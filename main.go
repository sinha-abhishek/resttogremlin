package main

import (
	"fmt"
	"html"
	"net/http"

	"bitbucket.org/abh_sinha/gremlin_client/gremlin"
)

func main() {
	client := gremlin.NewClient("localhost:8182")
	err := client.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	gr := client.NewGremlinRequest()
	gr.SetQuery("v = g.V().has(\"number\", '006').next() ; g.V(v.id()).inE('contact').property('has_app', false)")
	s, _ := client.SendRequest(gr)
	fmt.Println(s)
	// gr := client.NewGremlinRequest()
	// gr.SetQuery("g.V().has('uid',x)")
	// gr.AddBinding("x", 1004)
	// data, err2 := gr.PackageRequest()
	// fmt.Println(string(data))
	// if err2 != nil {
	// 	fmt.Println(err2)
	// 	return
	// }
	// s, _ := client.SendRequest(gr)
	// fmt.Println(s)
	// gr2 := client.NewGremlinRequest()
	// gr2.ReadQueryFromFileAndAddBindings("scripts/init.groovy", nil)
	// data2, err3 := gr2.PackageRequest()
	// fmt.Println(string(data2))
	// if err3 != nil {
	// 	fmt.Println(err3)
	// 	return
	// }
	// s2, _ := client.SendRequest(gr2)
	//
	// //s := client.SendRequest("{\"requestId\":\"655BD810-B41E-429D-B78F-3CC5F3B8E9BB\",\"processor\":\"\",\"op\":\"eval\",\"args\":{\"gremlin\":\"g.V().has('uid',1).values()\",\"language\":\"gremlin-groovy\"}}", "655BD810-B41E-429D-B78F-3CC5F3B8E9BA")
	// fmt.Println(s2)

	map1 := map[string]interface{}{"XNUM": "006", "XUID": 10060}
	gr4 := client.NewGremlinRequest()
	gr4.ReadQueryFromFileAndAddBindings("scripts/create_update_user_vertex.groovy", map1)
	s4, _ := client.SendRequest(gr4)
	fmt.Println(s4)
	// rmap := map[string]interface{}{"XNUM1": "004", "XNUM2": "006", "XNAME": "test6"}
	// gr3 := client.NewGremlinRequest()
	// gr3.ReadQueryFromFileAndAddBindings("scripts/add_contact", rmap)
	// s3, _ := client.SendRequest(gr3)
	// fmt.Println(s3)
	// rmap["XNUM2"] = "007"
	// rmap["XNAME"] = "test7"
	// gr3.ReadQueryFromFileAndAddBindings("scripts/add_contact", rmap)
	// s4, _ := client.SendRequest(gr3)
	// fmt.Println(s4)
	// router := mux.NewRouter().StrictSlash(true)
	// router.HandleFunc("/", Index)
	// log.Fatal(http.ListenAndServe(":8080", router))
}

/*
this is comment
*/
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}
