package resolver

import (
	"errors"
	"github.com/gin-gonic/gin"
	"reflect"
	"strconv"
)

type valueResolver struct {
	priority
	scope            Scope
	variable         string
	argumentPosition int
}

func Value(scope Scope, variable string, argumentPosition int) *valueResolver {
	return &valueResolver{
		priority:         priority{value: 256},
		scope:            scope,
		variable:         variable,
		argumentPosition: argumentPosition,
	}
}

// Path creates a value resolver which can inject a path variable into user handler argument.
func Path(variable string, argumentPosition int) *valueResolver {
	return &valueResolver{
		priority:         priority{value: 256},
		scope:            ScopePath,
		variable:         variable,
		argumentPosition: argumentPosition,
	}
}

// Query creates a value resolver which can inject a query variable into user handler argument.
func Query(variable string, argumentPosition int) *valueResolver {
	return &valueResolver{
		priority:         priority{value: 256},
		scope:            ScopeQuery,
		variable:         variable,
		argumentPosition: argumentPosition,
	}
}

// Header creates a value resolver which can inject a header variable into user handler argument.
func Header(variable string, argumentPosition int) *valueResolver {
	return &valueResolver{
		priority:         priority{value: 256},
		scope:            ScopeHeader,
		variable:         variable,
		argumentPosition: argumentPosition,
	}
}

func (r *valueResolver) CanResolve(ctx *gin.Context, argumentType reflect.Type, argument int) bool {
	if r.argumentPosition != argument {
		return false
	}

	switch r.scope {
	case ScopePath:
		_, exists := ctx.Params.Get(r.variable)
		return exists && isScalar(argumentType)
	case ScopeQuery:
		_, exits := ctx.GetQuery(r.variable)
		return exits && isScalar(argumentType)
	case ScopeHeader:
		return true
	}

	return false
}

func (r *valueResolver) Resolve(ctx *gin.Context, argumentType reflect.Type) (reflect.Value, error) {
	var val string
	switch r.scope {
	case ScopePath:
		val, _ = ctx.Params.Get(r.variable)
	case ScopeQuery:
		val, _ = ctx.GetQuery(r.variable)
	case ScopeHeader:
		val = ctx.GetHeader(r.variable)
	}

	// Convert value to the argumentPosition type.
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
