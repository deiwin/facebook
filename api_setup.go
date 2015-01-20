/*
Package facebook wraps the Facebook API in Go.

See the README on Github for more info: https://github.com/deiwin/facebook
*/
package facebook

import "net/http"

const (
	apiVersion = "v2.2"
)

type apiConf struct {
	graphURL string
}

// NewAPI creates an instance of the API using the provided *http.Client. It
// expects all the authentication to be handled by the http client.
func NewAPI(client *http.Client) API {
	return api{
		conf: apiConf{
			graphURL: "https://graph.facebook.com/" + apiVersion,
		},
		Client: client,
	}
}
