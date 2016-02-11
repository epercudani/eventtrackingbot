package requests

import (
	"github.com/kinslayere/eventtrackingbot/types"
	"github.com/kinslayere/eventtrackingbot/global"
	"errors"
	"fmt"
)

type GetUpdatesResponse struct {
	Ok			bool			`json:"ok"`
	Result		[]types.Update	`json:"result"`
}

type GetUpdatesRequest struct {
	getRequest *GetRequest
}

func NewGetUpdatesRequest() *GetUpdatesRequest {
	getRequest := NewGetRequest()
	getRequest.baseUrl = global.BASE_URL + "getUpdates"
	return &GetUpdatesRequest{getRequest}
}

func (r *GetUpdatesRequest) SetOffset(offset int64) {
	r.getRequest.SetParamInt64("offset", offset)
}

func (r *GetUpdatesRequest) SetTimeout(timeout int) {
	r.getRequest.SetParamInt("timeout", timeout)
}

func (r *GetUpdatesRequest) SetLimit(limit int) {
	r.getRequest.SetParamInt("limit", limit)
}

func (r *GetUpdatesRequest) Execute() (*GetUpdatesResponse, error) {

	var response GetUpdatesResponse
	err := r.getRequest.Execute(&response)
	if err != nil {
		return nil, err
	}

	if !response.Ok {
		return nil, errors.New(fmt.Sprintf("Error executing request '%v'", r.getRequest.baseUrl))
	}

	return &response, nil
}