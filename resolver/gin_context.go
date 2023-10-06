package resolver

import (
	"github.com/gin-gonic/gin"
	"reflect"
)

var ginContextType = reflect.TypeOf((*gin.Context)(nil))

type ginContextResolver struct {
	priority
}

// GinContext creates a resolver which can inject *gin.Context into user handler argumentPosition.
func GinContext() *ginContextResolver {
	return &ginContextResolver{priority{value: 200}}
}

func (r *ginContextResolver) CanResolve(_ *gin.Context, argumentType reflect.Type, _ int) bool {
	return argumentType == ginContextType
}

func (r *ginContextResolver) Resolve(ctx *gin.Context, _ reflect.Type) (reflect.Value, error) {
	return reflect.ValueOf(ctx), nil
}
