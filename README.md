# Facebook API wrapper in Go

[![GoDoc](https://godoc.org/github.com/deiwin/facebook?status.svg)](https://godoc.org/github.com/deiwin/facebook)

This is a work in progress. It is currently used by deiwin/luncher-api and only
implements the functionality needed by that project.

## Tutorial
### Configuration

Using `$FACEBOOK_APP_SECRET` and `$FACEBOOK_APP_ID` environment variables:
```go
redirectURL := "http://your.domain.com/login/facebook/redirected"
scopes := []string{"manage_pages", "publish_actions", "whatever"}
facebookConfig := facebook.NewConfig(redirectURL, scopes)
```

However, if you wish, you can also manage your app secret and id differently:
```go
facebookConfig := facebook.Config{
  AppID:       "your_app_id",
  AppSecret:   "your_app_secret",
  RedirectURL: "http://your.domain.com/login/facebook/redirected",
  Scopes:      []string{"manage_pages", "publish_actions", "whatever"},
}
```

### Authentication
#### Redirecting users to Facebook for login

```go
facebookAuthenticator := facebook.NewAuthenticator(facebookConfig)

// In an handler for the login request
session := "your_identifier_for_the_current_user_session"
redirectURL := facebookAuthenticator.AuthURL(session)
http.Redirect(w, r, redirectURL, http.StatusSeeOther)
```

#### Handling users redirected back from Facebook

```go
// In an handler for the RedirectURL in the configuration.
// E.g "http://your.domain.com/login/facebook/redirected"
session := "your_identifier_for_the_current_user_session"
tok, err := facebookAuthenticator.Token(session, r)
if err != nil {
  if err == facebook.ErrMissingState {
    http.Error(w, "Expecting a 'state' value", http.StatusBadRequest)
  } else if err == facebook.ErrInvalidState {
    http.Error(w, "Invalid 'state' value", http.StatusForbidden)
  } else if err == facebook.ErrMissingCode {
    http.Error(w, "Expecting a 'code' value", http.StatusBadRequest)
  }
  http.Error(w, "", http.StatusInternalServerError)
}
```

### Using the API
#### Receiving information about the current user

```go
api := fb.auth.APIConnection(tok)
user, err := api.Me()
if err != nil {
  ...
}
return user.ID
```

See the [GoDoc for the API](https://godoc.org/github.com/deiwin/facebook#API)
for more options.
