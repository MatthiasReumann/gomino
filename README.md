# gomino
Test-Utilities for [gin-gonic/gin](https://github.com/gin-gonic/gin).

## Examples
#### Most basic router from Gin's README.md
```golang 
func pingRouter(r *gin.Engine) *gin.Engine {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	return r
}

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
```

#### Using middlewares
```golang
func userRouter(r *gin.Engine) *gin.Engine {
	r.GET("/user", func(c *gin.Context) {
		if c.MustGet("session-username").(string) == "hansi" {
			c.JSON(http.StatusOK, gin.H{
				"message": "hello hansi",
			})
		} else {
			c.AbortWithStatus(http.StatusForbidden)
		}
	})
	return r
}

func TestWithMiddleware(t *testing.T) {
	testCases := gomino.TestCases{
		"user hansi": {
			Method: http.MethodGet,
			Url:    "/user",
			Middlewares: []func(c *gin.Context){
				func(c *gin.Context) {
					c.Set("session-username", "hansi")
				},
			},
			ExpectedCode:     http.StatusOK,
			ExpectedResponse: gin.H{"message": "hello hansi"},
		},
		"user not hansi": {
			Method: http.MethodGet,
			Url:    "/user",
			Middlewares: []func(c *gin.Context){
				func(c *gin.Context) {
					c.Set("session-username", "bobby")
				},
			},
			ExpectedCode: http.StatusForbidden,
		},
	}
	testCases.Run(t, func(r *gin.Engine) {
		userRouter(r)
	})
}
```
