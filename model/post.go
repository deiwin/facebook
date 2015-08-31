package model

import (
	"fmt"
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
	Post
	IsPublished bool `json:"is_published,omitempty"`
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
		form.Set("scheduled_publish_time", p.ScheduledPublishTime.Format(time.RFC3339))
	}
	if (p.BackdatedTime != time.Time{}) {
		form.Set("backdated_time", p.BackdatedTime.Format(time.RFC3339))
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
