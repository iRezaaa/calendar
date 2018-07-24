package model

import "gopkg.in/mgo.v2/bson"

type Note struct {
	ID     bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	Title  string        `json:"title,omitempty" bson:"title,omitempty"`
	Body   string        `json:"body,omitempty" bson:"body,omitempty"`
	UserID string        `json:"-" bson:"user_id,omitempty"`
	Day    int           `json:"-,omitempty" bson:"day,omitempty"`
	Month  int           `json:"-,omitempty" bson:"month,omitempty"`
	Year   int           `json:"-,omitempty" bson:"year,omitempty"`
}
