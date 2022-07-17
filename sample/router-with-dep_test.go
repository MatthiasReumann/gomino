package sample

import (
	"github.com/gin-gonic/gin"
	"gomino/gomino"
	"net/http"
	"testing"
)

type UserDao interface {
	Get() string
}

type userDao struct {
	name string
}

func (u userDao) Get() string {
	return u.name
}

func NewUserDaoMock(name string) UserDao {
	return userDao{name}
}

func loginRouter(r *gin.Engine, dao UserDao) *gin.Engine {
	r.POST("/login", func(c *gin.Context) {
		if dao.Get() == "hansi" {
			c.Status(http.StatusOK)
		} else {
			c.AbortWithStatus(http.StatusForbidden)
		}
	})
	return r
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
