package facebook

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/deiwin/facebook/model"
)

// API provides access to the Facebook API graph methods
type API interface {
	// GET /me
	//
	// https://developers.facebook.com/docs/graph-api/reference/v2.3/user#read
	Me() (*model.User, error)

	// GET /me/accounts
	//
	// https://developers.facebook.com/docs/graph-api/reference/v2.3/user/accounts#read
	Accounts() (*model.Accounts, error)

	// GET /{page-id}
	//
	// https://developers.facebook.com/docs/graph-api/reference/page/
	Page(pageID string) (*model.Page, error)

	// POST /{page-id}/feed
	//
	// https://developers.facebook.com/docs/graph-api/reference/v2.3/page/feed#publish
	PagePublish(pageAccessToken, pageID, message string) (*model.Post, error)

	// DELETE /{post-id}
	//
	// https://developers.facebook.com/docs/graph-api/reference/v2.3/post#deleting
	PostDelete(pageAccessToken, postID string) error
}

type api struct {
	*http.Client
	conf apiConf
}

func (a api) Me() (*model.User, error) {
	resp, err := a.get("/me", url.Values{})
	if err != nil {
		return nil, err
	}
	var user model.User
	err = json.Unmarshal(resp, &user)
	return &user, err
}

func (a api) Accounts() (*model.Accounts, error) {
	resp, err := a.get("/me/accounts", url.Values{})
	if err != nil {
		return nil, err
	}
	var accs model.Accounts
	err = json.Unmarshal(resp, &accs)
	return &accs, err
}

func (a api) Page(pageID string) (*model.Page, error) {
	resp, err := a.get("/"+pageID, url.Values{})
	if err != nil {
		return nil, err
	}
	var page model.Page
	err = json.Unmarshal(resp, &page)
	return &page, err
}

func (a api) PagePublish(pageAccessToken, pageID, message string) (*model.Post, error) {
	resp, err := a.post(fmt.Sprintf("/%s/feed", pageID), url.Values{
		"message":      {message},
		"access_token": {pageAccessToken},
		// TODO add publish time
	})
	if err != nil {
		return nil, err
	}
	var post model.Post
	err = json.Unmarshal(resp, &post)
	return &post, err
}

func (a api) PostDelete(pageAccessToken, postID string) error {
	return a.delete(postID, url.Values{
		"access_token": {pageAccessToken},
	})
}

func (a api) delete(path string, data url.Values) error {
	url := fmt.Sprintf("%s/%s?%s", a.conf.graphURL, path, data.Encode())
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	if _, err = a.Do(req); err != nil {
		return err
	}
	return nil
}

func (a api) get(path string, data url.Values) ([]byte, error) {
	url := fmt.Sprintf("%s/%s?%s", a.conf.graphURL, path, data.Encode())
	resp, err := a.Get(url)
	return parseResponse(resp, err)
}

func (a api) post(path string, data url.Values) ([]byte, error) {
	resp, err := a.PostForm(a.conf.graphURL+path, data)
	return parseResponse(resp, err)
}

func parseResponse(resp *http.Response, err error) ([]byte, error) {
	if err != nil {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
		return nil, err
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return bytes, err
	}
	return bytes, parseError(bytes)
}

func parseError(bytes []byte) error {
	var resp struct {
		Error *Error `json:"error,omitempty"`
	}
	err := json.Unmarshal(bytes, &resp)
	if err != nil {
		return err
	}
	return resp.Error
}

type Error struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    int    `json:"code"`
}

func (e *Error) Error() string {
	return e.Message
}
