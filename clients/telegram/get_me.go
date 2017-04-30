package telegram

import (
	"github.com/kinslayere/eventtrackingbot/global"
	"github.com/kinslayere/eventtrackingbot/clients"
	"errors"
	"fmt"

)

type GetMeResult struct {
	Id		int			`json:"id"`
	FirstName	string			`json:"first_name"`
	UserName	string			`json:"username"`
}

type GetMeResponse struct {
	Ok		bool			`json:"ok"`
	Result		GetMeResult		`json:"result"`
}

type GetMeRequest struct {
	getRequest *requests.GetRequest
}

func NewGetMeRequest() *GetMeRequest {
	getRequest := requests.NewGetRequest()
	getRequest.SetBaseURL(global.TELEGRAM_BASE_URL + "getMe")

	return &GetMeRequest{getRequest}
}

func (r *GetMeRequest) Execute() (*GetMeResponse, error) {

	var response GetMeResponse
	err := r.getRequest.Execute(&response)
	if err != nil {
		return nil, err
	}

	if !response.Ok {
		return nil, errors.New(fmt.Sprintf("Error executing request '%v'", r.getRequest.GetBaseURL()))
	}

	return &response, nil
}