package types

type Venue struct {
	id		string		`json:"id"`
	name		string		`json:"name"`
	contact		Contact		`json:"contact"`
	location	Location	`json:"location"`
}

type Contact struct {
	twitter 	string		`json:"twitter"`
	phone		string		`json:"phone"`
	formattedPhone	string		`json:"formattedPhone"`
}

type Location struct {
	address		string		`json:"address"`
	crossStreet	string		`json:"crossStreet"`
	city		string		`json:"city"`
	state		string		`json:"state"`
	postalCode	string		`json:"postalCode"`
	country		string		`json:"country"`
	lat		float64		`json:"lat"`
	lng		float64		`json:"lng"`
	distance	float64		`json:"distance"`
}