package sample

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func pingRouter(r *gin.Engine) *gin.Engine {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	return r
}

func main() {
	r := gin.Default()
	pingRouter(r).Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
