# gomino
![Go Report Card](https://goreportcard.com/badge/github.com/matthiasreumann/gomino)

Gomino provides test-utilities for [gin-gonic/gin](https://github.com/gin-gonic/gin)'s web framework.

## Examples
#### Simple 'ping' router from Gin's README.md
```golang 
func pingRouter(r *gin.Engine) *gin.Engine {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
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
	}.Run(t)
}
```

#### Using middlewares
```golang
func userRouter(r *gin.Engine) {
	r.GET("/user", func(c *gin.Context) {
		if c.MustGet("session-username").(string) == "hansi" {
			c.JSON(http.StatusOK, gin.H{"message": "hello hansi"})
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
            ExpectedResponse: gin.H{"message": "hello hansi"},
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
```

#### Using router-functions with dependency injection

```golang
func loginRouter(r *gin.Engine, dao UserDao) {
	r.POST("/login", func(c *gin.Context) {
		if dao.Get() == "hansi" {
			c.Status(http.StatusOK)
		} else {
			c.AbortWithStatus(http.StatusForbidden)
		}
	})
}

func TestRouterWithDependencies(t *testing.T) {
	gomino.TestCases{
		"user hansi": {
			Router: func(r *gin.Engine) {
				loginRouter(r, NewUserDaoMock("hansi"))
			},
			Method:       http.MethodPost,
			Url:          "/login",
			ExpectedCode: http.StatusOK,
		},
		"user not hansi": {
			Router: func(r *gin.Engine) {
				loginRouter(r, NewUserDaoMock("bobby"))
			},
			Method:       http.MethodPost,
			Url:          "/login",
			ExpectedCode: http.StatusForbidden,
		},
	}.Run(t)
}
```
