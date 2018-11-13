package model

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type News struct {
	ID             bson.ObjectId   `json:"_id,omitempty" bson:"_id,omitempty"`
	Title          string          `json:"title,omitempty" bson:"title,omitempty"`
	Content        string          `json:"content,omitempty" bson:"content,omitempty"`
	ImageURL       string          `json:"image_url,omitempty" bson:"image_url,omitempty"`
	CreateTime     time.Time       `json:"create_time,omitempty" bson:"create_time,omitempty"`
	LastUpdateTime time.Time       `json:"last_update_time,omitempty" bson:"last_update_time,omitempty"`
}