package foursquare

import (
	"github.com/kinslayere/eventtrackingbot/global"
	"github.com/kinslayere/eventtrackingbot/clients"
	"errors"
	"fmt"
)

type VenueSearchResponse struct {
	Ok		bool		`json:"ok"`
	//Result		[]foursquare.Update	`json:"result"`
}

type VenueSearchRequest struct {
	getRequest *requests.GetRequest
}

func NewVenueSearchRequest() *VenueSearchRequest {
	getRequest := requests.NewGetRequest()
	getRequest.SetBaseURL(global.FOURSQUARE_BASE_URL + "venues/search")

	vsr := VenueSearchRequest{getRequest}

	vsr.SetClientId(global.FOURSQUARE_CLIENT_ID)
	vsr.SetClientSecret(global.FOURSQUARE_CLIENT_SECRET)
	vsr.SetMode(global.FOURSQUARE_MODE)
	vsr.SetVersion(global.FOURSQUARE_VERSION)

	return &vsr
}

func (r *VenueSearchRequest) SetClientId(clientId string) {
	r.getRequest.SetParamString("client_id", clientId)
}

func (r *VenueSearchRequest) SetClientSecret(clientSecret string) {
	r.getRequest.SetParamString("client_secret", clientSecret)
}

func (r *VenueSearchRequest) SetVersion(version string) {
	r.getRequest.SetParamString("v", version)
}

func (r *VenueSearchRequest) SetMode(mode string) {
	r.getRequest.SetParamString("m", mode)
}

func (r *VenueSearchRequest) Execute() (*VenueSearchResponse, error) {

	var response VenueSearchResponse
	err := r.getRequest.Execute(&response)
	if err != nil {
		return nil, err
	}

	if !response.Ok {
		return nil, errors.New(fmt.Sprintf("Error executing request '%v'", r.getRequest.GetBaseURL()))
	}

	return &response, nil
}