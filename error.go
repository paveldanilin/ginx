package ginx

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

type ErrorInterceptor interface {
	Intercept(e ErrorEvent)
}

type ErrorEvent interface {
	Context() *gin.Context
	Error() error
	Response() Response
}

type errorEvent struct {
	ctx *gin.Context
	err error
	res Response
}

func newErrorEvent(ctx *gin.Context, err any) ErrorEvent {
	return &errorEvent{
		ctx: ctx,
		err: valueToError(err),
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

func valueToError(v any) error {
	switch t := v.(type) {
	case error:
		return t
	case string:
		return errors.New(t)
	default:
		return fmt.Errorf("%v", v)
	}
}
