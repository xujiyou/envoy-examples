package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.GET("/authz/*action", func(c *gin.Context) {
		authzHeader := c.Request.Header.Get("Authorization")
		if authzHeader == "my-token" {
			c.Writer.WriteHeader(200)
		} else {
			c.Writer.WriteHeader(401)
		}
	})
	_ = r.Run("0.0.0.0:6060")
}
