package examples

import (
	"github.com/gin-gonic/gin"
	"gomino"
	"net/http"
	"testing"
)

func userRouter(r *gin.Engine) {
	r.GET("/user", func(c *gin.Context) {
		if c.MustGet("session-username").(string) == "hansi" {
			c.JSON(http.StatusOK, gin.H{
				"message": "hello hansi!",
			})
		} else {
			c.AbortWithStatus(http.StatusForbidden)
		}
	})
}

func TestWithMiddleware(t *testing.T) {
	gomino.TestCases{
		"user hansi": {
			Router: userRouter,
			Method: http.MethodGet,
			Url:    "/user",
			Middlewares: []func(c *gin.Context){
				func(c *gin.Context) {
					c.Set("session-username", "hansi")
				},
			},
			ExpectedCode:     http.StatusOK,
			ExpectedResponse: gin.H{"message": "hello hansi!"},
		},
		"user not hansi": {
			Router: userRouter,
			Method: http.MethodGet,
			Url:    "/user",
			Middlewares: []func(c *gin.Context){
				func(c *gin.Context) {
					c.Set("session-username", "bobby")
				},
			},
			ExpectedCode: http.StatusForbidden,
		},
	}.Run(t)
}
