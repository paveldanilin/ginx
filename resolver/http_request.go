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

// HttpRequest creates a resolver which can inject *http.Request into user handler argumentPosition.
func HttpRequest() *httpRequestResolver {
	return &httpRequestResolver{priority{value: 200}}
}

func (r *httpRequestResolver) CanResolve(_ *gin.Context, argumentType reflect.Type, _ int) bool {
	return argumentType == httpRequestType
}

func (r *httpRequestResolver) Resolve(ctx *gin.Context, _ reflect.Type) (reflect.Value, error) {
	return reflect.ValueOf(ctx.Request), nil
}
