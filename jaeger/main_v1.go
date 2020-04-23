package main

import (
	"github.com/emicklei/go-restful"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

// This example shows the minimal code needed to get a restful.WebService working.
//
// GET http://localhost:8081/hello

func main() {
	ws := new(restful.WebService)
	ws.Route(ws.GET("/hello").To(hello))
	restful.Add(ws)
	log.Println("server start in 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func hello(req *restful.Request, resp *restful.Response) {
	log.Println("request hello v1")
	bodyResp, _ := http.Get("http://service-v2:9001/hello")
	body, _ := ioutil.ReadAll(bodyResp.Body)
	_, _ = io.WriteString(resp, "world v1"+string(body[:]))
}
