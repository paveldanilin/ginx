package ginx

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

type ErrorInterceptor interface {
	InterceptError(Error)
}

type ErrorInterceptorFunc func(Error)

func (f ErrorInterceptorFunc) InterceptError(e Error) {
	f(e)
}

type Error interface {
	Context() *gin.Context
	Error() error
	Response() Response
}

type errorEvent struct {
	ctx *gin.Context
	err error
	res Response
}

func newError(ctx *gin.Context, err any) *errorEvent {
	return &errorEvent{
		ctx: ctx,
		err: anyToError(err),
		res: NewResponse(500),
	}
}

func (e errorEvent) Context() *gin.Context {
	return e.ctx
}

func (e errorEvent) Error() error {
	return e.err
}

func (e errorEvent) Response() Response {
	return e.res
}

func anyToError(v any) error {
	switch t := v.(type) {
	case error:
		return t
	case string:
		return errors.New(t)
	default:
		return fmt.Errorf("%v", v)
	}
}
