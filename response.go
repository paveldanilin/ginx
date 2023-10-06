package ginx

import "reflect"

type Response interface {
	Status() int
	SetStatus(status int)
	ContentType() string
	SetContentType(contentType string)
	Body() any
	SetBody(body any)
}

type response struct {
	status      int
	contentType string
	body        any
}

func NewResponse(status int) Response {
	return &response{
		status:      status,
		contentType: "",
		body:        nil,
	}
}

func OKResponse() Response {
	return &response{
		status: 200,
	}
}

func (r *response) Status() int {
	return r.status
}

func (r *response) SetStatus(status int) {
	if status == 0 {
		status = 200
	}
	r.status = status
}

func (r *response) ContentType() string {
	return r.contentType
}

func (r *response) SetContentType(contentType string) {
	r.contentType = contentType
}

func (r *response) Body() any {
	return r.body
}

func (r *response) SetBody(body any) {
	r.body = body
}

func (response) fromValue(v reflect.Value) Response {
	res := NewResponse(200)

	switch v.Type().Kind() {
	// Number is a status code
	case reflect.Int:
		res.SetStatus(int(v.Int()))
		if res.Body() == nil {
			res.SetBody("")
		}
	case reflect.Uint:
		res.SetStatus(int(v.Uint()))
		if res.Body() == nil {
			res.SetBody("")
		}
	case reflect.String:
		res.SetBody(v.String())
	case reflect.Pointer:
		if v.IsNil() {
			res.SetBody(nil)
		} else {
			return response{}.fromValue(reflect.Indirect(v))
		}
	case reflect.Array, reflect.Slice:
		if v.Len() == 0 {
			res.SetBody([]any{})
		} else {
			res.SetBody(v.Interface())
		}
	case reflect.Interface, reflect.Struct:
		if r, isResponse := v.Interface().(Response); isResponse {
			res.SetStatus(r.Status())
			res.SetContentType(r.ContentType())
			res.SetBody(r.Body())
		} else {
			res.SetBody(v.Interface())
		}
	}

	return res
}
