package requests

import (
	"log"
	"errors"
	"bytes"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"strconv"
	"net/url"
)

type Request struct {
	url string
	params map[string]interface{}
}

func NewRequest() *Request {
	var r Request
	r.params = make(map[string]interface{})
	return &r
}

func (r *Request) SetUrl(url string) {
	r.url = url
}

func (r *Request) AddParamInt64(name string, param int64) {
	r.params[name] = strconv.FormatInt(param, 10)
}

func (r *Request) AddParamString(name, param string) {
	r.params[name] = url.QueryEscape(param)
}

func (r *Request) AddParamBoolean(name string, param bool) {
	r.params[name] = string(param)
}

func (r *Request) AddParamStringer(name string, param fmt.Stringer) {
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

func (r *Request) DoRequest(response *interface{}) (err error) {

	if _, ok := r.params["chat_id"]; !ok {
		return errors.New("chat_id is required")
	}

	if _, ok := r.params["text"]; !ok {
		return errors.New("text is required")
	}

	paramsString := r.GetParamsString()
	url := r.url + "?" + paramsString
	log.Printf("Request:url: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Request:response:body: %s", body)

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatal(err)
	}

	return
}