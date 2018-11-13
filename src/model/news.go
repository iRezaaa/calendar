package model

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type NewsContentType int

const (
	NewsContentTypePlain NewsContentType = 0x0001
	NewsContentTypeHTML  NewsContentType = 0x0002
)

type News struct {
	ID             bson.ObjectId   `json:"_id,omitempty" bson:"_id,omitempty"`
	Title          string          `json:"title,omitempty" bson:"title,omitempty"`
	Content        string          `json:"content,omitempty" bson:"content,omitempty"`
	ContentType    NewsContentType `json:"content_type,omitempty" bson:"content_type,omitempty"`
	ImageURL       string          `json:"image_url,omitempty" bson:"image_url,omitempty"`
	CreateTime     time.Time       `json:"create_time,omitempty" bson:"create_time,omitempty"`
	LastUpdateTime time.Time       `json:"last_update_time,omitempty" bson:"last_update_time,omitempty"`
}