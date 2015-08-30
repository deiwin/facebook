package model

// Accounts is used to unmarshal the /me/accounts response
//
// https://developers.facebook.com/docs/graph-api/reference/v2.4/user/accounts#fields
type Accounts struct {
	Data []Page `json:"data"`
}

// Page represents a Facebook Page
//
// https://developers.facebook.com/docs/graph-api/reference/v2.4/page#readfields
type Page struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Location    Location `json:"location,omitempty"`
	Phone       string   `json:"phone,omitempty"`
	Website     string   `json:"website,omitempty"`
	Emails      []string `json:"emails,omitempty"`
	AccessToken string   `json:"access_token,omitempty"`
}

// Location holds the location information for a Facebook object, including the address
// and the geographical location.
//
// https://developers.facebook.com/docs/graph-api/reference/location/
type Location struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	Country string `json:"country"`
}
