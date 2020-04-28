package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.GET("/v1/certs/list/approved", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"certificates": []map[string]string{
				{
					"fingerprint_sha256": "a45e46487977c62cd60a3a4e8ec044f8fe16115bd1c8f5cddf3d99f82dc864a7",
				},
			},
		})
	})
	_ = r.Run("0.0.0.0:6060")
}
