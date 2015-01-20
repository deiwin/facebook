package facebook

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/deiwin/facebook/model"
)

// API provides access to the Facebook API graph methods
type API interface {
	// /me
	//
	// https://developers.facebook.com/docs/graph-api/reference/v2.2/user#read
	Me() (model.User, error)
	// /me/accounts
	//
	// https://developers.facebook.com/docs/graph-api/reference/v2.2/user/accounts#read
	Accounts() (model.Accounts, error)
}

type api struct {
	*http.Client
	conf apiConf
}

func (a api) Me() (user model.User, err error) {
	resp, err := a.get("/me")
	if err != nil {
		return
	}
	err = json.Unmarshal(resp, &user)
	return
}

func (a api) Accounts() (accs model.Accounts, err error) {
	resp, err := a.get("/me/accounts")
	if err != nil {
		return
	}
	err = json.Unmarshal(resp, &accs)
	return
}

func (a api) get(path string) (response []byte, err error) {
	resp, err := a.Get(a.conf.graphURL + path)
	if err != nil {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
		return
	}
	defer resp.Body.Close()
	response, err = ioutil.ReadAll(resp.Body)
	return
}
