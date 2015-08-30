package model

import "time"

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

// Link is used as a pointer to images in a post
//
// https://developers.facebook.com/docs/graph-api/reference/v2.4/link
type Link struct {
	Link string `json:"link"`
}
