package model

// Post represents a Facebook Post
//
// https://developers.facebook.com/docs/graph-api/reference/v2.2/post
type Post struct {
	ID      string `json:"id"`
	Message string `json:"message,omitempty"`
}
