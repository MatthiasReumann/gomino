package gomino

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// TestCases shall contain all test cases of a single test suite, e.g. for one particular endpoint.
// The key represents the test name, the value an instance of TestCase
type TestCases map[string]TestCase

// HttpHeader is a key value map for HTTP header fields such as Content-Type, Cache-Control,...
type HttpHeader map[string]string

// TestCase contains everything a single test needs to execute
type TestCase struct {
	Router      func(*gin.Engine)
	Method      string
	Url         string
	Middlewares []func(c *gin.Context)

	ContentType string
	Body        interface{}

	ExpectedHeader   HttpHeader
	ExpectedCode     int
	ExpectedResponse interface{}

	Before func()
	After  func()
}

// Equal is the function signiture for one's favourite testing framework
type Equal func(*testing.T, interface{}, interface{})

// Run executes all tests of a given TestCases object
func (tc TestCases) Run(t *testing.T, equal Equal) {
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

			testCase.Router(r)

			c.Request, _ = http.NewRequest(testCase.Method, testCase.Url, testCase.getBody())
			if testCase.Method != http.MethodGet {
				c.Request.Header.Set("Content-Type", testCase.getContentType())
			}
			r.ServeHTTP(w, c.Request)

			for field, expected := range testCase.ExpectedHeader {
				equal(t, expected, w.Header().Get(field))
			}

			equal(t, testCase.ExpectedCode, w.Code)

			if testCase.ExpectedResponse != nil {
				equal(t, testCase.getResponse(), w.Body.Bytes())
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
