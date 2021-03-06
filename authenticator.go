package facebook

import (
	"errors"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

var (
	ErrMissingState = errors.New("A Facebook redirect request is missing the 'state' value")
	ErrInvalidState = errors.New("A Facebook redirect request's 'state' value does not match the session")
	ErrMissingCode  = errors.New("A Facebook redirect request is missing the 'code' value")
	ErrNoSuchPage   = errors.New("The user does not have access to that page")
)

// Authenticator provides the authentication functionality for Facebook users
// using Facebook's OAuth
type Authenticator interface {
	// AuthURL returns a Facebook URL the user should be redirect to. The user
	// will then be asked to log in by Facebook at that URL and will be redirected
	// back to the configured RedirectURL.
	AuthURL(state string) string
	// Token get's the longer term user access token from the redirect request.
	// Also checks that the provided state matches that of the redirect request and
	// returns "", ErrInvalidState if it doesn't.
	Token(state string, r *http.Request) (*oauth2.Token, error)
	// APIConnection returns an API instance that can be used to make authenticated
	// requests to the Facebook API.
	APIConnection(tok *oauth2.Token) API
	// PageAccessToken retrieves a page access token of the specified page if the
	// user has access to that page and returns "", ErrNoSuchPage if they don't.
	PageAccessToken(tok *oauth2.Token, pageID string) (string, error)
}

// NewAuthenticator initializes and returns an Authenticator
func NewAuthenticator(conf Config) Authenticator {
	opts := &oauth2.Config{
		ClientID:     conf.AppID,
		ClientSecret: conf.AppSecret,
		RedirectURL:  conf.RedirectURL,
		Scopes:       conf.Scopes,
		Endpoint:     facebook.Endpoint,
	}
	return authenticator{
		Config: opts,
	}
}

type authenticator struct {
	*oauth2.Config
}

func (a authenticator) AuthURL(state string) string {
	return a.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (a authenticator) Token(state string, r *http.Request) (*oauth2.Token, error) {
	expectedState := r.FormValue("state")
	if expectedState == "" {
		return nil, ErrMissingState
	} else if expectedState != state {
		return nil, ErrInvalidState
	}
	code := r.FormValue("code")
	if code == "" {
		return nil, ErrMissingCode
	}
	return a.Exchange(oauth2.NoContext, code)
}

func (a authenticator) PageAccessToken(tok *oauth2.Token, pageID string) (string, error) {
	connection := a.APIConnection(tok)
	accs, err := connection.Accounts()
	if err != nil {
		return "", err
	}
	for _, page := range accs.Data {
		if page.ID == pageID {
			return page.AccessToken, nil
		}
	}
	return "", ErrNoSuchPage
}

func (a authenticator) APIConnection(tok *oauth2.Token) API {
	client := a.Config.Client(oauth2.NoContext, tok)
	return NewAPI(client)
}
