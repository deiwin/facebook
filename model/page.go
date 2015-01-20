package model

// Accounts is used to unmarshal the /me/accounts response
type Accounts struct {
	Data []Page `json:"data"`
}

// Page represents a Facebook Page
//
// https://developers.facebook.com/docs/graph-api/reference/v2.2/page#readfields
type Page struct {
	ID string `json:"id"`
	// Additional field included in the /me/accounts response
	//
	// https://developers.facebook.com/docs/graph-api/reference/v2.2/user/accounts#fields
	AccessToken string `json:"access_token,omitempty"`
}
