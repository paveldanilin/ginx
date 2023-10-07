package resolver

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
)

var httpRequestType = reflect.TypeOf((*http.Request)(nil))

type httpRequestResolver struct {
	priority
}

// HttpRequest resolver injects *http.Request into user handler.
//
//	controller.GET("/foo", func(req *http.Request) {
//		// Do something with http.Request
//	})
func HttpRequest() *httpRequestResolver {
	return &httpRequestResolver{priority{value: 200}}
}

func (r *httpRequestResolver) CanResolve(_ *gin.Context, argumentType reflect.Type, _ int) bool {
	return argumentType == httpRequestType
}

func (r *httpRequestResolver) Resolve(ctx *gin.Context, _ reflect.Type) (reflect.Value, error) {
	return reflect.ValueOf(ctx.Request), nil
}
