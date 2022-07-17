package sample

import (
	"github.com/gin-gonic/gin"
	"gomino/gomino"
	"net/http"
	"testing"
)

func pingRouter(r *gin.Engine) *gin.Engine {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	return r
}

func TestGinReadme(t *testing.T) {
	router := func(r *gin.Engine) {
		pingRouter(r)
	}
	gomino.TestCases{
		"ping": {
			Router:           router,
			Method:           http.MethodGet,
			Url:              "/ping",
			ExpectedCode:     http.StatusOK,
			ExpectedResponse: gin.H{"message": "pong"},
		},
	}.Run(t)
}
