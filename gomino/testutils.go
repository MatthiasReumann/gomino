package gomino

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type TestCases map[string]TestCase

type TestCase struct {
	Method           string
	Url              string
	Middlewares      []func(c *gin.Context)
	ContentType      string
	Body             interface{}
	ExpectedCode     int
	ExpectedResponse interface{}

	Before func()
	After  func()
}

func (tc TestCases) Run(t *testing.T, router func(r *gin.Engine)) {
	for name, testCase := range tc {
		if testCase.Before != nil {
			testCase.Before()
		}

		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, r := gin.CreateTestContext(w)

			for _, mw := range testCase.Middlewares {
				r.Use(mw)
			}

			router(r)

			c.Request, _ = http.NewRequest(testCase.Method, testCase.Url, testCase.getBody())
			if testCase.Method != http.MethodGet {
				c.Request.Header.Set("Content-Type", testCase.getContentType())
			}
			r.ServeHTTP(w, c.Request)

			assert.Equal(t, testCase.ExpectedCode, w.Code)

			if testCase.ExpectedResponse != nil {
				fmt.Println(string(testCase.getResponse()))
				fmt.Println(w.Body.String())
				assert.Equal(t, testCase.getResponse(), w.Body.Bytes())
			}
		})

		if testCase.After != nil {
			testCase.After()
		}
	}
}

func (c TestCase) getContentType() string {
	if len(c.ContentType) > 0 {
		return c.ContentType
	}

	return "application/json"
}

func (c TestCase) getBody() io.Reader {
	if c.Body == nil {
		return bytes.NewBuffer([]byte{})
	}
	switch c.Body.(type) {
	case io.Reader:
		return c.Body.(io.Reader)
	case string:
		return bytes.NewBufferString(c.Body.(string))
	case []byte:
		return bytes.NewBuffer(c.Body.([]byte))
	default:
		j, err := json.Marshal(c.Body)
		if err != nil {
			panic(errors.New("invalid body type"))
		}

		return bytes.NewBuffer(j)
	}
}

func (c TestCase) getResponse() []byte {
	switch c.ExpectedResponse.(type) {
	case string:
		return c.ExpectedResponse.([]byte)
	case []byte:
		return c.ExpectedResponse.([]byte)
	default:
		j, err := json.Marshal(c.ExpectedResponse)
		if err != nil {
			panic(errors.New("invalid expectedResponse type"))
		}

		return j
	}
}

func NewMultipartFormData(fieldName, fileName string) (bytes.Buffer, *multipart.Writer) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	file, _ := os.Open(fileName)
	fw, _ := w.CreateFormFile(fieldName, file.Name())
	io.Copy(fw, file)
	w.Close()
	return b, w
}

func First(a interface{}, b interface{}) interface{} {
	return a
}

func Second(a interface{}, b interface{}) interface{} {
	return b
}
