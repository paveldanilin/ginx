package resolver

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/textproto"
	"reflect"
	"strconv"
)

type valueResolver struct {
	priority
	scope            Scope
	variable         string
	argumentPosition int
	defaultValue     any
}

func Value(scope Scope, variable string, argumentPosition int, defaultValue any) *valueResolver {
	return &valueResolver{
		priority:         priority{value: 250},
		scope:            scope,
		variable:         variable,
		argumentPosition: argumentPosition,
		defaultValue:     defaultValue,
	}
}

// Path creates a resolver which can inject a path value.
func Path(pathVariable string, argumentPosition int) *valueResolver {
	return Value(ScopePath, pathVariable, argumentPosition, nil)
}

// Query creates a resolver which can inject a query value.
func Query(queryVariable string, argumentPosition int) *valueResolver {
	return Value(ScopeQuery, queryVariable, argumentPosition, nil)
}

func QueryOrDefault(queryVariable string, argumentPosition int, defaultValue any) *valueResolver {
	return Value(ScopeQuery, queryVariable, argumentPosition, defaultValue)
}

// Header creates a resolver which can inject a header value.
func Header(headerVariable string, argumentPosition int) *valueResolver {
	return Value(ScopeHeader, headerVariable, argumentPosition, nil)
}

func (r *valueResolver) CanResolve(_ *gin.Context, argumentType reflect.Type, argument int) bool {
	return r.argumentPosition == argument && isScalar(argumentType)
}

func (r *valueResolver) defaultValueToArgumentType(defaultValue any, argumentType reflect.Type) reflect.Value {
	if defaultValue == nil {
		switch argumentType.Kind() {
		case reflect.String:
			return reflect.ValueOf("")
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return reflect.ValueOf(0)
		case reflect.Float32, reflect.Float64:
			return reflect.ValueOf(0.0)
		case reflect.Bool:
			return reflect.ValueOf(false)
		case reflect.Pointer:
			return reflect.ValueOf(nil)
		}
	}
	return reflect.ValueOf(defaultValue)
}

func (r *valueResolver) Resolve(ctx *gin.Context, argumentType reflect.Type) (reflect.Value, error) {
	var val string
	switch r.scope {
	case ScopePath:
		if v, exists := ctx.Params.Get(r.variable); exists {
			val = v
		} else {
			return r.defaultValueToArgumentType(r.defaultValue, argumentType), nil
		}
	case ScopeQuery:
		if v, exists := ctx.GetQuery(r.variable); exists {
			val = v
		} else {
			return r.defaultValueToArgumentType(r.defaultValue, argumentType), nil
		}
	case ScopeHeader:
		if _, exists := ctx.Request.Header[textproto.CanonicalMIMEHeaderKey(r.variable)]; exists {
			val = ctx.GetHeader(r.variable)
		} else {
			return r.defaultValueToArgumentType(r.defaultValue, argumentType), nil
		}
	}

	// Convert value to the argument type.
	switch argumentType.Kind() {
	case reflect.String:
		return reflect.ValueOf(val), nil
	case reflect.Bool:
		b, err := strconv.ParseBool(val)
		if err != nil {
			return reflect.ValueOf(false), err
		}
		return reflect.ValueOf(b), nil
	case reflect.Int:
		i, err := strconv.Atoi(val)
		if err != nil {
			return reflect.ValueOf(0), err
		}
		return reflect.ValueOf(i), nil
	case reflect.Int8:
		i, err := strconv.ParseInt(val, 10, 8)
		if err != nil {
			return reflect.ValueOf(0), err
		}
		return reflect.ValueOf(i), nil
	case reflect.Int16:
		i, err := strconv.ParseInt(val, 10, 16)
		if err != nil {
			return reflect.ValueOf(0), err
		}
		return reflect.ValueOf(i), nil
	case reflect.Int32:
		i, err := strconv.ParseInt(val, 10, 32)
		if err != nil {
			return reflect.ValueOf(0), err
		}
		return reflect.ValueOf(i), nil
	case reflect.Int64:
		i, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return reflect.ValueOf(0), err
		}
		return reflect.ValueOf(i), nil
	case reflect.Float32:
		f, err := strconv.ParseFloat(val, 32)
		if err != nil {
			return reflect.ValueOf(0.0), err
		}
		return reflect.ValueOf(f), nil
	case reflect.Float64:
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return reflect.ValueOf(0.0), err
		}
		return reflect.ValueOf(f), nil
	case reflect.Uint:
		i, err := strconv.ParseUint(val, 10, 32)
		if err != nil {
			return reflect.ValueOf(0), err
		}
		return reflect.ValueOf(i), nil
	case reflect.Uint8:
		i, err := strconv.ParseUint(val, 10, 8)
		if err != nil {
			return reflect.ValueOf(0), err
		}
		return reflect.ValueOf(i), nil
	case reflect.Uint16:
		i, err := strconv.ParseUint(val, 10, 16)
		if err != nil {
			return reflect.ValueOf(0), err
		}
		return reflect.ValueOf(i), nil
	case reflect.Uint32:
		i, err := strconv.ParseUint(val, 10, 32)
		if err != nil {
			return reflect.ValueOf(0), err
		}
		return reflect.ValueOf(i), nil
	case reflect.Uint64:
		i, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return reflect.ValueOf(0), err
		}
		return reflect.ValueOf(i), nil
	}

	return reflect.ValueOf(nil), errors.New("could not resolve")
}

func isScalar(t reflect.Type) bool {
	return t.Kind() == reflect.Int ||
		t.Kind() == reflect.Int8 ||
		t.Kind() == reflect.Int16 ||
		t.Kind() == reflect.Int32 ||
		t.Kind() == reflect.Int64 ||
		t.Kind() == reflect.Float32 ||
		t.Kind() == reflect.Float64 ||
		t.Kind() == reflect.Uint ||
		t.Kind() == reflect.Uint8 ||
		t.Kind() == reflect.Uint16 ||
		t.Kind() == reflect.Uint32 ||
		t.Kind() == reflect.Uint64 ||
		t.Kind() == reflect.String ||
		t.Kind() == reflect.Bool
}
