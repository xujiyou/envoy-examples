package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	_ = r.RunTLS("0.0.0.0:9090", "/home/admin/k8s-cluster/envoy/ssl/cert/server.crt", "/home/admin/k8s-cluster/envoy/ssl/cert/server.key")
}
