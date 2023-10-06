package ginx

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/paveldanilin/ginx/resolver"
	"github.com/paveldanilin/ginx/slices"
	"reflect"
	"runtime"
	"strings"
)

type Controller struct {
	BasePath          string
	ContentType       string
	router            *gin.Engine
	handlerMap        map[string]*handler
	argumentResolvers []ArgumentResolver
	middlewares       []gin.HandlerFunc
	tester            *Tester
	errorInterceptor  ErrorInterceptor
}

func NewController(r *gin.Engine) *Controller {
	return &Controller{
		router:            r,
		handlerMap:        map[string]*handler{},
		argumentResolvers: []ArgumentResolver{},
		middlewares:       []gin.HandlerFunc{},
		tester:            NewTester(r),
	}
}

func NewDefaultController(r *gin.Engine) *Controller {
	c := &Controller{
		router:            r,
		handlerMap:        map[string]*handler{},
		argumentResolvers: []ArgumentResolver{},
		middlewares:       []gin.HandlerFunc{},
		tester:            NewTester(r),
	}

	// HttpRequest creates a resolver which can inject *http.Request into user handler argument.
	c.Use(resolver.HttpRequest())
	// GinContext creates a resolver which can inject *gin.Context into user handler argument.
	c.Use(resolver.GinContext())
	// Context creates a resolver which can inject gin.Context.Request.Context() into user handler argument.
	c.Use(resolver.Context())
	// Bind request body into struct, bind query/header/path values by tag 'ginx'.
	c.Use(resolver.StructResolver())

	return c
}

func (c *Controller) Use(opt HandlerOption) {
	if r, isResolver := opt.(ArgumentResolver); isResolver {
		c.argumentResolvers = append(c.argumentResolvers, r)
		return
	}

	if middleware, isMiddleware := opt.(gin.HandlerFunc); isMiddleware {
		c.middlewares = append(c.middlewares, middleware)
		return
	}
}

func (c *Controller) GET(path string, handler HandlerFunc, opts ...HandlerOption) error {
	return c.registerHandler("GET", path, handler, opts...)
}

func (c *Controller) POST(path string, handler HandlerFunc, opts ...HandlerOption) error {
	return c.registerHandler("POST", path, handler, opts...)
}

func (c *Controller) PUT(path string, handler HandlerFunc, opts ...HandlerOption) error {
	return c.registerHandler("PUT", path, handler, opts...)
}

func (c *Controller) PATCH(path string, handler HandlerFunc, opts ...HandlerOption) error {
	return c.registerHandler("PATCH", path, handler, opts...)
}

func (c *Controller) DELETE(path string, handler HandlerFunc, opts ...HandlerOption) error {
	return c.registerHandler("DELETE", path, handler, opts...)
}

func (c *Controller) HEAD(path string, handler HandlerFunc, opts ...HandlerOption) error {
	return c.registerHandler("HEAD", path, handler, opts...)
}

func (c *Controller) OPTIONS(path string, handler HandlerFunc, opts ...HandlerOption) error {
	return c.registerHandler("OPTIONS", path, handler, opts...)
}

func (c *Controller) Tester() *Tester {
	return c.tester
}

func (c *Controller) registerHandler(method, path string, handlerFunc HandlerFunc, opts ...HandlerOption) error {
	handlerFuncReflect := reflect.ValueOf(handlerFunc)
	if handlerFuncReflect.Kind() != reflect.Func {
		return errors.New("handler must be function")
	}

	method = normalizeHttpMethod(method)
	path = normalizePath(path)

	h := &handler{
		name:      getGoMethodName(handlerFuncReflect.Pointer()),
		numIn:     handlerFuncReflect.Type().NumIn(),
		numOut:    handlerFuncReflect.Type().NumOut(),
		function:  handlerFuncReflect,
		arguments: []reflect.Type{},
		resolvers: []ArgumentResolver{},
	}

	h.init(c.argumentResolvers, opts...)

	c.handlerMap[getHandlerId(method, path)] = h

	// [<controller.middlewares>, <handler.middlewares>, <request.handler>]
	var ginHandlers []gin.HandlerFunc = slices.Join(c.middlewares, handlerOptions(opts).Middlewares())
	ginHandlers = append(ginHandlers, c.handleRequest)

	fullPath := c.BasePath + path

	switch method {
	case "GET":
		c.router.GET(fullPath, ginHandlers...)
	case "POST":
		c.router.POST(fullPath, ginHandlers...)
	case "PUT":
		c.router.PUT(fullPath, ginHandlers...)
	case "PATCH":
		c.router.PATCH(fullPath, ginHandlers...)
	case "DELETE":
		c.router.DELETE(fullPath, ginHandlers...)
	case "HEAD":
		c.router.HEAD(fullPath, ginHandlers...)
	case "OPTIONS":
		c.router.OPTIONS(fullPath, ginHandlers...)
	}

	return errors.New("unknown method")
}

func (c *Controller) handleRequest(ctx *gin.Context) {
	h := c.getHandler(ctx)
	if h == nil {
		panic("handler not found")
	}

	ctx.Set("ginx_handler_name", h.name)
	ctx.Set("ginx_controller_response_type", c.ContentType)
	ctx.Set("ginx_handler_response_type", h.responseContentType)

	hArgs, err := h.resolveArguments(ctx)
	if err != nil {
		panic(err)
	}

	// TODO: validate arguments

	handlerResponse := h.function.Call(hArgs)

	c.response(ctx, handlerResponse)
}

func (c *Controller) response(ctx *gin.Context, handlerResponse []reflect.Value) {
	switch len(handlerResponse) {
	case 0:
		// If user handler returns void
		c.sendResponse(ctx, OKResponse())
	case 1:
		// If user handler returns: (<error>)
		if isError(handlerResponse[0]) {
			e := newErrorEvent(ctx, handlerResponse[0])
			if c.errorInterceptor != nil {
				c.errorInterceptor.Intercept(e)
			}
			c.sendResponse(ctx, e.Response())
			return
		}

		// TODO: handle return <int> ?

		// If user handler returns: (<userdata>)
		c.sendResponse(ctx, response{}.fromValue(handlerResponse[0]))
	default:
		// If user handler returns: (<userdata>, <error>)
		if isError(handlerResponse[1]) {
			e := newErrorEvent(ctx, handlerResponse[1].Interface())
			if c.errorInterceptor != nil {
				c.errorInterceptor.Intercept(e)
			}
			c.sendResponse(ctx, e.Response())
			return
		}
		// If user handler returns: (<error>, <userdata>)
		if isError(handlerResponse[0]) {
			e := newErrorEvent(ctx, handlerResponse[0].Interface())
			if c.errorInterceptor != nil {
				c.errorInterceptor.Intercept(e)
			}
			c.sendResponse(ctx, e.Response())
			return
		}

		res := response{}.fromValue(handlerResponse[0])

		// If user handler returns: (<userdata>, <int|uint>)
		// <userdata>, <http response status>
		if status, isStatus := getStatus(handlerResponse[1]); isStatus {
			res.SetStatus(status)
		}

		c.sendResponse(ctx, res)
	}
}

func (c *Controller) sendResponse(ctx *gin.Context, res Response) {
	responseContentType := c.getResponseContentType(ctx, res)

	if s, isString := res.Body().(string); isString {
		ctx.Data(res.Status(), responseContentType, []byte(s))
		return
	}

	if b, isBytes := res.Body().([]byte); isBytes {
		ctx.Data(res.Status(), responseContentType, b)
		return
	}

	switch getFormat(responseContentType) {
	case "json":
		ctx.JSON(res.Status(), res.Body())
	case "xml":
		ctx.XML(res.Status(), res.Body())
	case "text":
		ctx.Data(res.Status(), responseContentType, []byte(fmt.Sprintf("%v", res.Body())))
	default:
		ctx.Data(res.Status(), responseContentType, []byte(fmt.Sprintf("%v", res.Body())))
	}
}

func (c *Controller) getResponseContentType(ctx *gin.Context, res Response) string {
	if strings.TrimSpace(res.ContentType()) != "" {
		return res.ContentType()
	}

	handlerResponseContentType, existsHandlerContentType := ctx.Get("ginx_handler_response_type")
	if existsHandlerContentType && strings.TrimSpace(handlerResponseContentType.(string)) != "" {
		return handlerResponseContentType.(string)
	}

	controllerResponseContentType, existsControllerContentType := ctx.Get("ginx_controller_response_type")
	if existsControllerContentType && strings.TrimSpace(controllerResponseContentType.(string)) != "" {
		return controllerResponseContentType.(string)
	}

	// TODO: get content type by request header

	return gin.MIMEPlain
}

func (c *Controller) getHandler(ctx *gin.Context) *handler {
	if h, exists := c.handlerMap[getHandlerId(ctx.Request.Method, ctx.FullPath())]; exists {
		return h
	}
	return nil
}

func getGoMethodName(m uintptr) string {
	return strings.TrimSuffix(runtime.FuncForPC(m).Name(), "-fm")
}

func normalizeHttpMethod(method string) string {
	return strings.ToUpper(strings.TrimSpace(method))
}

func normalizePath(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		path = "/"
	}
	return path
}

func getHandlerId(method, path string) string {
	return method + path
}

func getFormat(contentType string) string {
	if strings.Contains(contentType, "json") {
		return "json"
	}
	if strings.Contains(contentType, "xml") {
		return "xml"
	}
	if strings.Contains(contentType, "text") {
		return "text"
	}
	return "text"
}

var errType = reflect.TypeOf((*error)(nil)).Elem()

func isError(v reflect.Value) bool {
	if v.Interface() == nil {
		return false
	}
	return v.Type().Implements(errType)
}

func getStatus(v reflect.Value) (int, bool) {
	if v.Kind() == reflect.Int ||
		v.Kind() == reflect.Int16 ||
		v.Kind() == reflect.Int32 ||
		v.Kind() == reflect.Int64 ||
		v.Kind() == reflect.Uint ||
		v.Kind() == reflect.Uint16 ||
		v.Kind() == reflect.Uint32 ||
		v.Kind() == reflect.Uint64 {
		return int(v.Int()), true
	}
	return 0, false
}
