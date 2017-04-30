package requests

import (
	"strconv"
	"net/url"
	"fmt"
	"bytes"
)

type Request struct {
	baseUrl string
	params  map[string]string
}

func NewRequest() *Request {
	var r Request
	r.params = make(map[string]string)
	return &r
}

func (r *Request) SetBaseURL(value string) {
	r.baseUrl = value
}

func (r *Request) GetBaseURL() string {
	return r.baseUrl
}

func (r *Request) SetParamInt(name string, param int) {
	r.params[name] = strconv.Itoa(param)
}

func (r *Request) SetParamInt64(name string, param int64) {
	r.params[name] = strconv.FormatInt(param, 10)
}

func (r *Request) SetParamString(name, param string) {
	r.params[name] = url.QueryEscape(param)
}

func (r *Request) SetParamBoolean(name string, param bool) {
	r.params[name] = strconv.FormatBool(param)
}

func (r *Request) SetParamStringer(name string, param fmt.Stringer) {
	r.params[name] = url.QueryEscape(param.String())
}

func (r *Request) GetParamsString() string {

	var paramsString bytes.Buffer
	first := true

	for k, v := range r.params {
		if first {
			paramsString.WriteString(fmt.Sprintf("%s=%v", k, v))
			first = false
		} else {
			paramsString.WriteString(fmt.Sprintf("&%s=%v", k, v))
		}
	}

	return paramsString.String()
}

func (r *Request) HasParam(name string) bool {
	_, ok := r.params[name]
	return ok
}