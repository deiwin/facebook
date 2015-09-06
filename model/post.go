package model

import (
	"fmt"
	"io"
	"net/url"
	"strconv"
	"time"
)

// Post represents a Facebook Post
//
// https://developers.facebook.com/docs/graph-api/reference/v2.4/post
type Post struct {
	ID                   string    `json:"id"`
	Message              string    `json:"message,omitempty"`
	Published            bool      `json:"published"`
	ScheduledPublishTime time.Time `json:"scheduled_publish_time,omitempty"`
	BackdatedTime        time.Time `json:"backdated_time,omitempty"`
	ObjectAttachment     string    `json:"object_attachment,omitempty"`
	ChildAttachments     []Link    `json:"child_attachments,omitempty"`
}

// PostResponse also holds fields specific only to the Post object when
// it's received as a response from Facebook, unlike the Post object which is
// used to post to Facebook.
type PostResponse struct {
	ID                   string    `json:"id"`
	Message              string    `json:"message"`
	IsPublished          bool      `json:"is_published"`
	ScheduledPublishTime time.Time `json:"scheduled_publish_time"`
	CreatedTime          time.Time `json:"created_time"`
}

type PhotoResponse struct {
	ID     string `json:"id"`
	PostID string `json:"post_id"`
}

type Photo struct {
	Post
	Photo io.Reader
}

// Link is used as a pointer to images in a post
//
// https://developers.facebook.com/docs/graph-api/reference/v2.4/link
type Link struct {
	Link string `json:"link"`
}

func (p *Post) AsForm() url.Values {
	form := url.Values{}
	if p.Message != "" {
		form.Set("message", p.Message)
	}
	form.Set("published", strconv.FormatBool(p.Published))
	if (p.ScheduledPublishTime != time.Time{}) {
		form.Set("scheduled_publish_time", strconv.FormatInt(p.ScheduledPublishTime.Unix(), 10))
	}
	if (p.BackdatedTime != time.Time{}) {
		form.Set("backdated_time", strconv.FormatInt(p.BackdatedTime.Unix(), 10))
	}
	if p.ObjectAttachment != "" {
		form.Set("object_attachment", p.ObjectAttachment)
	}
	for i, childAttachment := range p.ChildAttachments {
		// This isn't well documented, but is the correct format for an array of objects
		form.Set(fmt.Sprintf("%s[%d][%s]", "child_attachments", i, "link"), childAttachment.Link)
	}
	return form
}

// func (p *Photo) AsForm() url.Values {
// 	return p.Post.AsForm()
// }
