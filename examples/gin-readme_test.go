package examples

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/matthiasreumann/gomino"
	"net/http"
	"testing"
)

func pingRouter(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
}

func TestGinReadme(t *testing.T) {
	gomino.TestCases{
		"ping": {
			Router:           pingRouter,
			Method:           http.MethodGet,
			Url:              "/ping",
			ExpectedCode:     http.StatusOK,
			ExpectedResponse: gin.H{"message": "pong"},
		},
	}.Run(t, assert.Equal)
}
