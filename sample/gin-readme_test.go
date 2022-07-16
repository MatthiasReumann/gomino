package sample

import (
	"github.com/gin-gonic/gin"
	"gomino/gomino"
	"net/http"
	"testing"
)

func TestGinReadme(t *testing.T) {
	testCases := gomino.TestCases{
		"ping": {
			Method:           http.MethodGet,
			Url:              "/ping",
			ExpectedCode:     http.StatusOK,
			ExpectedResponse: gin.H{"message": "pong"},
		},
	}
	testCases.Run(t, func(r *gin.Engine) {
		pingRouter(r)
	})
}
