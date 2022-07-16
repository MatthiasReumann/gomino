# gomino
Test-Utilities for [gin-gonic/gin](https://github.com/gin-gonic/gin).

## Examples
```golang 
// gin-gonic/gin README web router
func pingRouter(r *gin.Engine) *gin.Engine {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	return r
}
...
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
