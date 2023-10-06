package resolver

import (
	"context"
	"github.com/gin-gonic/gin"
	"reflect"
)

var contextType = reflect.TypeOf((*context.Context)(nil)).Elem()

type contextResolver struct {
	priority
}

// Context creates a resolver which can inject gin.Context.Request.Context() into user handler argumentPosition.
func Context() *contextResolver {
	return &contextResolver{priority{value: 200}}
}

func (r *contextResolver) CanResolve(_ *gin.Context, argumentType reflect.Type, _ int) bool {
	return argumentType == contextType
}

func (r *contextResolver) Resolve(ctx *gin.Context, _ reflect.Type) (reflect.Value, error) {
	return reflect.ValueOf(ctx.Request.Context()), nil
}
