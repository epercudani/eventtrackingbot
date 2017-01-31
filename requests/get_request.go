package requests

import (
	"log"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

type GetRequest struct {
	*Request
}

func NewGetRequest() *GetRequest {
	r := NewRequest()
	return &GetRequest{r}
}

func (r *GetRequest) Execute(response interface{}) (err error) {

	url := r.baseUrl
	paramsString := r.GetParamsString()
	if len(paramsString) > 0 {
		url += "?" + paramsString
	}

	//log.Printf("Request:url: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	//log.Printf("Request:response:body: %s", body)

	return json.Unmarshal(body, &response)
}