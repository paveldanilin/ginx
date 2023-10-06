package resolver

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/paveldanilin/ginx/requestbody"
	"reflect"
	"strconv"
	"strings"
)

const tagKey = "ginx"

var requestBodyType = reflect.TypeOf((*requestbody.RequestBody)(nil)).Elem()
var empty = reflect.ValueOf(nil)

type tagParam struct {
	name  string
	value string
}

type structResolver struct {
	priority
}

func StructResolver() *structResolver {
	return &structResolver{priority{value: 200}}
}

func (r structResolver) CanResolve(_ *gin.Context, argumentType reflect.Type, _ int) bool {
	return argumentType.Kind() == reflect.Struct || argumentType.Kind() == reflect.Map
}

func (r structResolver) Resolve(ctx *gin.Context, argumentType reflect.Type) (reflect.Value, error) {
	arg := reflect.New(argumentType)
	indirect := reflect.Indirect(arg)

	// If HTTP method support body (PUT, POST, PATCH) and a target type implements 'requestbody.RequestBody' interface.
	err := r.bindBody(ctx, arg)
	if err != nil {
		return empty, err
	}

	if argumentType.Kind() == reflect.Struct {
		r.processTag(ctx, indirect)
	}

	return indirect, nil
}

func (r structResolver) processTag(ctx *gin.Context, val reflect.Value) {
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)

		// Nested
		if field.Type.Kind() == reflect.Struct {
			r.processTag(ctx, val.Field(i))
		}

		tagDef, hasTag := field.Tag.Lookup(tagKey)
		if !hasTag {
			continue
		}

		fieldValue := val.FieldByName(field.Name)
		if !fieldValue.IsValid() || !fieldValue.CanSet() {
			continue
		}

		tagParams := r.parseTag(tagDef)
		for _, tParam := range tagParams {
			switch tParam.name {
			case string(ScopeQuery):
				v := ctx.Query(tParam.value)
				r.setField(fieldValue, v)
			case string(ScopePath):
				v := ctx.Param(tParam.value)
				r.setField(fieldValue, v)
			case string(ScopeHeader):
				v := ctx.GetHeader(tParam.value)
				r.setField(fieldValue, v)
			}
		}
	}
}

func (r structResolver) parseTag(tagDef string) []tagParam {
	tagEntries := strings.Split(tagDef, ",")

	var params []tagParam

	for i := 0; i < len(tagEntries); i++ {
		kv := strings.SplitN(tagEntries[i], "=", 2)
		if len(kv) != 2 {
			continue
		}
		params = append(params, tagParam{name: kv[0], value: kv[1]})
	}

	return params
}

func (r structResolver) setField(v reflect.Value, s string) {
	switch v.Kind() {
	case reflect.String:
		v.SetString(s)
	case reflect.Bool:
		if s == "" {
			v.SetBool(false)
		} else {
			b, _ := strconv.ParseBool(s)
			v.SetBool(b)
		}
	case reflect.Int, reflect.Int32, reflect.Int64:
		if s == "" {
			v.SetInt(0)
		} else {
			i, _ := strconv.Atoi(s)
			v.SetInt(int64(i))
		}
	case reflect.Float32:
		if s == "" {
			v.SetFloat(0)
		} else {
			f, _ := strconv.ParseFloat(s, 32)
			v.SetFloat(f)
		}
	case reflect.Float64:
		if s == "" {
			v.SetFloat(0)
		} else {
			f, _ := strconv.ParseFloat(s, 64)
			v.SetFloat(f)
		}
	}
}

func (r structResolver) bindBody(ctx *gin.Context, out reflect.Value) error {
	if !httpMethodHasBody(ctx.Request.Method) || !out.Type().Implements(requestBodyType) {
		return nil
	}

	bodyFormat := out.Elem().MethodByName("RequestBodyFormat").Call([]reflect.Value{})[0].String()

	if bodyFormat == "" {
		bodyFormat = r.getFormat(ctx)
	}

	switch bodyFormat {
	case "json":
		err := ctx.BindJSON(out.Interface())
		if err != nil {
			return err
		}
		return nil
	case "xml":
		err := ctx.BindXML(out.Interface())
		if err != nil {
			return err
		}
		return nil
	}

	return fmt.Errorf("unnown body format '%s'", bodyFormat)
}

func (r structResolver) getFormat(ctx *gin.Context) string {
	if strings.Contains(ctx.ContentType(), "json") {
		return "json"
	}

	if strings.Contains(ctx.ContentType(), "xml") {
		return "xml"
	}

	return ctx.ContentType()
}

func httpMethodHasBody(method string) bool {
	return method == "POST" || method == "PUT" || method == "PATCH"
}
