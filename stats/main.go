package main

import (
	"github.com/emicklei/go-restful"
	"io"
	"log"
	"net/http"
)

// This example shows the minimal code needed to get a restful.WebService working.
//
// GET http://localhost:8082/hello

func main() {
	ws := new(restful.WebService)
	ws.Route(ws.GET("/hello").To(hello))
	restful.Add(ws)
	log.Println("server start in 8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}

func hello(req *restful.Request, resp *restful.Response) {
	log.Println("request hello")
	_, _ = io.WriteString(resp, "world")
}
