package model

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Notification struct {
	ID         bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	Title      string        `json:"title,omitempty" bson:"title,omitempty"`
	Body       string        `json:"body,omitempty" bson:"body,omitempty"`
	SendTime   time.Time     `json:"send_time,omitempty" bson:"send_time,omitempty"`
	StatusCode int           `json:"status_code,omitempty" bson:"status_code,omitempty"`
	PushID     string        `json:"-" bson:"push_id"`
}
