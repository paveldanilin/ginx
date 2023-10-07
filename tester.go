package ginx

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
)

type Tester struct {
	router *gin.Engine
}

func NewTester(r *gin.Engine) *Tester {
	return &Tester{router: r}
}

func (t *Tester) POST(url string, body []byte) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", url, bytes.NewReader(body))
	t.router.ServeHTTP(w, req)
	return w
}

func (t *Tester) POSTJson(url string, body any) *httptest.ResponseRecorder {
	bodyData, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", url, bytes.NewReader(bodyData))
	t.router.ServeHTTP(w, req)
	return w
}

func (t *Tester) PUT(url string, body []byte) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", url, bytes.NewReader(body))
	t.router.ServeHTTP(w, req)
	return w
}

func (t *Tester) PATCH(url string, body []byte) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", url, bytes.NewReader(body))
	t.router.ServeHTTP(w, req)
	return w
}

func (t *Tester) DELETE(url string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", url, nil)
	t.router.ServeHTTP(w, req)
	return w
}

func (t *Tester) GET(url string, headers map[string]string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", url, nil)

	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	t.router.ServeHTTP(w, req)
	return w
}
