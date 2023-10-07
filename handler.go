package ginx

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/paveldanilin/ginx/slices"
	"reflect"
	"sort"
)

// ArgumentResolver represents a request handler argument resolver.
type ArgumentResolver interface {
	// Priority returns a resolver priority.
	Priority() int

	// CanResolve checks whether resolver can resolve value or not.
	CanResolve(*gin.Context, reflect.Type, int) bool

	// Resolve resolves value.
	Resolve(*gin.Context, reflect.Type) (reflect.Value, error)
}

// HandlerFunc represents a request handler.
type HandlerFunc any

type HandlerOption any

func Produce(contentType string) func(*handler) {
	return func(h *handler) {
		h.responseContentType = contentType
	}
}

func ProduceJSON() func(*handler) {
	return func(h *handler) {
		h.responseContentType = gin.MIMEJSON
	}
}

func ProduceXML() func(*handler) {
	return func(h *handler) {
		h.responseContentType = gin.MIMEXML
	}
}

type handlerOptions []HandlerOption

func (o handlerOptions) Middlewares() []gin.HandlerFunc {
	m := slices.Filter(o, func(_ int, t HandlerOption) bool {
		_, isMiddleware := t.(gin.HandlerFunc)
		return isMiddleware
	})

	return slices.Map(m, func(t HandlerOption) gin.HandlerFunc {
		return t.(gin.HandlerFunc)
	})
}

// handler represents a request handler definition.
type handler struct {
	name                string
	responseContentType string
	numIn               int
	numOut              int
	function            reflect.Value
	arguments           []reflect.Type
	resolvers           []ArgumentResolver
}

func (h *handler) init(controllerArgumentResolvers []ArgumentResolver, opts ...HandlerOption) {
	// Fill arguments
	for i := 0; i < h.numIn; i++ {
		h.arguments = append(h.arguments, h.function.Type().In(i))
	}

	// Bind argument resolvers
	h.resolvers = append(h.resolvers, controllerArgumentResolvers...)

	for _, opt := range opts {
		if resolver, isResolver := opt.(ArgumentResolver); isResolver {
			h.resolvers = append(h.resolvers, resolver)
		} else if optFunc, isFunc := opt.(func(*handler)); isFunc {
			optFunc(h)
		}
	}

	// Sort resolvers by priority
	sort.Slice(h.resolvers, func(i, j int) bool {
		return h.resolvers[i].Priority() > h.resolvers[j].Priority()
	})
}

func (h *handler) findArgumentResolver(ctx *gin.Context, argumentType reflect.Type, argumentIndex int) ArgumentResolver {
	resolver, resolverPresent := slices.First(h.resolvers, func(t ArgumentResolver) bool {
		return t.CanResolve(ctx, argumentType, argumentIndex)
	})
	if resolverPresent {
		return resolver
	}
	return nil
}

func (h *handler) resolveArguments(ctx *gin.Context) ([]reflect.Value, error) {
	var args []reflect.Value

	for i, at := range h.arguments {
		argumentPosition := i + 1
		argumentType := at

		// TODO: can we cache this?
		resolver := h.findArgumentResolver(ctx, argumentType, argumentPosition)
		if resolver == nil {
			return nil, fmt.Errorf("resolver not found for argument at position [%d]", argumentPosition)
		}

		resolved, err := resolver.Resolve(ctx, argumentType)
		if err != nil {
			return nil, err
		}

		args = append(args, resolved)
	}

	return args, nil
}
