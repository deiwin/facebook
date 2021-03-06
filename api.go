package facebook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/deiwin/facebook/model"
)

const ISO8601 = "2006-01-02T15:04:05-0700"

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
	PagePublish(pageAccessToken, pageID string, post *model.Post) (*model.PostResponse, error)

	// POST /{page_id}/photos
	//
	// https://developers.facebook.com/docs/graph-api/reference/page/photos/#Creating
	PagePhotoCreate(pageAccessToken, pageID string, photo *model.Photo) (*model.PhotoResponse, error)

	// GET /{post-id}
	//
	// https://developers.facebook.com/docs/graph-api/reference/v2.4/post#read
	Post(pageAccessToken, postID string) (*model.PostResponse, error)

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

func (a api) PagePublish(pageAccessToken, pageID string, post *model.Post) (*model.PostResponse, error) {
	resp, err := a.postFormable(fmt.Sprintf("/%s/feed", pageID), url.Values{
		"access_token": {pageAccessToken},
	}, post)
	if err != nil {
		return nil, err
	}
	var respPost model.PostResponse
	err = json.Unmarshal(resp, &respPost)
	return &respPost, err
}

func (a api) PagePhotoCreate(pageAccessToken, pageID string, photo *model.Photo) (*model.PhotoResponse, error) {
	additionalData := photo.AsForm()
	additionalData.Add("access_token", pageAccessToken)
	resp, err := a.postMultipartForm(fmt.Sprintf("/%s/photos", pageID), "photo.jpeg", photo.Photo, additionalData)
	if err != nil {
		return nil, err
	}
	var respPhoto model.PhotoResponse
	err = json.Unmarshal(resp, &respPhoto)
	return &respPhoto, err
}

func (a api) Post(pageAccessToken, postID string) (*model.PostResponse, error) {
	fieldsToInclude := getJsonTagsForType(reflect.TypeOf(model.PostResponse{}))
	resp, err := a.getFields("/"+postID, url.Values{
		"access_token": {pageAccessToken},
	}, fieldsToInclude)
	if err != nil {
		return nil, err
	}
	// Wrap the actual PostResponse object to do necessary conversions
	var post struct {
		model.PostResponse
		ScheduledPublishTime int64  `json:"scheduled_publish_time"`
		CreatedTime          string `json:"created_time"`
	}
	err = json.Unmarshal(resp, &post)
	postResponse := post.PostResponse
	if post.ScheduledPublishTime != 0 {
		postResponse.ScheduledPublishTime = time.Unix(post.ScheduledPublishTime, 0)
	}
	if post.CreatedTime != "" {
		if postResponse.CreatedTime, err = time.Parse(ISO8601, post.CreatedTime); err != nil {
			return nil, err
		}
	}
	return &postResponse, err
}

func (a api) PostUpdate(pageAccessToken, postID string, post *model.Post) error {
	_, err := a.postFormable("/"+postID, url.Values{
		"access_token": {pageAccessToken},
	}, post)
	return err
}

func (a api) PostDelete(pageAccessToken, postID string) error {
	return a.delete(postID, url.Values{
		"access_token": {pageAccessToken},
	})
}

func getJsonTagsForType(t reflect.Type) []string {
	tags := make([]string, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		tags[i] = t.Field(i).Tag.Get("json")
	}
	return tags
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

func (a api) getFields(path string, data url.Values, fields []string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s?%s&fields=%s", a.conf.graphURL, path, data.Encode(), strings.Join(fields, ","))
	resp, err := a.Get(url)
	return parseResponse(resp, err)
}

func (a api) post(path string, data url.Values) ([]byte, error) {
	resp, err := a.PostForm(a.conf.graphURL+path, data)
	return parseResponse(resp, err)
}

func (a api) postMultipartForm(path string, filename string, reader io.Reader, additionalData url.Values) ([]byte, error) {
	multipartReader, contentType, err := createMultipartForm(filename, reader, additionalData)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", a.conf.graphURL+path, multipartReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)

	resp, err := a.Do(req)
	return parseResponse(resp, err)
}

func createMultipartForm(filename string, reader io.Reader, additionalData url.Values) (*bytes.Buffer, string, error) {
	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	multipartWriter := multipart.NewWriter(&b)
	defer multipartWriter.Close()

	formFileWriter, err := multipartWriter.CreateFormFile("source", filename)
	if err != nil {
		return nil, "", err
	}
	if _, err = io.Copy(formFileWriter, reader); err != nil {
		return nil, "", err
	}
	for key, vals := range additionalData {
		for _, val := range vals {
			formFieldWriter, err := multipartWriter.CreateFormField(key)
			if err != nil {
				return nil, "", err
			}
			if _, err = formFieldWriter.Write([]byte(val)); err != nil {
				return nil, "", err
			}
		}
	}
	return &b, multipartWriter.FormDataContentType(), nil
}

func (a api) postFormable(path string, data url.Values, formable model.Formable) ([]byte, error) {
	form := formable.AsForm()
	for k, vs := range data {
		for _, v := range vs {
			form.Add(k, v)
		}
	}
	resp, err := a.PostForm(a.conf.graphURL+path, form)
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
