package facebook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"

	"github.com/deiwin/facebook/model"
)

// API provides access to the Facebook API graph methods
type API interface {
	// GET /me
	//
	// https://developers.facebook.com/docs/graph-api/reference/v2.4/user#read
	Me() (*model.User, error)

	// GET /me/accounts
	//
	// https://developers.facebook.com/docs/graph-api/reference/v2.4/user/accounts#read
	Accounts() (*model.Accounts, error)

	// GET /{page-id}
	//
	// https://developers.facebook.com/docs/graph-api/reference/page/
	Page(pageID string) (*model.Page, error)

	// POST /{page-id}/feed
	//
	// https://developers.facebook.com/docs/graph-api/reference/v2.4/page/feed#publish
	PagePublish(pageAccessToken, pageID string, post *model.Post) (*model.Post, error)

	// GET /{post-id}
	//
	// https://developers.facebook.com/docs/graph-api/reference/v2.4/post#read
	Post(pageAccessToken, postID string) (*model.Post, error)

	// POST /{post-id}
	//
	// https://developers.facebook.com/docs/graph-api/reference/v2.4/post#updating
	PostUpdate(pageAccessToken, postID string, post *model.Post) error

	// DELETE /{post-id}
	//
	// https://developers.facebook.com/docs/graph-api/reference/v2.4/post#deleting
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

func (a api) PagePublish(pageAccessToken, pageID string, post *model.Post) (*model.Post, error) {
	resp, err := a.postJSON(fmt.Sprintf("/%s/feed", pageID), url.Values{
		"access_token": {pageAccessToken},
	}, post)
	if err != nil {
		return nil, err
	}
	var respPost model.Post
	err = json.Unmarshal(resp, &respPost)
	return &respPost, err
}

func (a api) Post(pageAccessToken, postID string) (*model.Post, error) {
	resp, err := a.get("/"+postID, url.Values{
		"access_token": {pageAccessToken},
	})
	if err != nil {
		return nil, err
	}
	var post model.Post
	err = json.Unmarshal(resp, &post)
	return &post, err
}

func (a api) PostUpdate(pageAccessToken, postID string, post *model.Post) error {
	_, err := a.postJSON("/"+postID, url.Values{
		"access_token": {pageAccessToken},
	}, post)
	return err
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

func (a api) postJSON(path string, params url.Values, v interface{}) ([]byte, error) {
	url := fmt.Sprintf("%s/%s?%s", a.conf.graphURL, path, params.Encode())
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	resp, err := a.Client.Post(url, "application/json", bytes.NewReader(data))
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
	// If we were to just return resp.Error, then in case an error didn't actually occur, the returned
	// value would be `(*facebook.Error)(nil)`, which doesn't equal to nil and will make it look like
	// there's always an error
	if !reflect.ValueOf(resp.Error).IsNil() {
		return resp.Error
	}
	return nil
}

type Error struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    int    `json:"code"`
}

func (e *Error) Error() string {
	return e.Message
}
