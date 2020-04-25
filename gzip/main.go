package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message1": "pong",
			"message2": "pong",
			"message3": "pong",
			"message4": "pong",
			"message5": "pong",
			"message6": "pong",
		})
	})
	_ = r.Run("0.0.0.0:9090")
}
