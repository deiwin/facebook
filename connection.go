package facebook

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/deiwin/luncher-api/facebook/model"
)

// Connection provides access to the Facebook API graph methods
type Connection interface {
	// /me
	Me() (model.User, error)
	// /me/accounts
	Accounts() (model.Accounts, error)
}

type connection struct {
	*http.Client
	api api
}

func (c connection) Me() (user model.User, err error) {
	resp, err := c.get("/me")
	if err != nil {
		return
	}
	err = json.Unmarshal(resp, &user)
	return
}

func (c connection) Accounts() (accs model.Accounts, err error) {
	resp, err := c.get("/me/accounts")
	if err != nil {
		return
	}
	err = json.Unmarshal(resp, &accs)
	return
}

func (c connection) get(path string) (response []byte, err error) {
	resp, err := c.Get(c.api.graphURL + path)
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
