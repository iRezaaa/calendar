package model

import "gopkg.in/mgo.v2/bson"

type Banner struct {
	ID       bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	ImageURL string        `json:"image_url,omitempty" bson:"image_url,omitempty"`
	Title    string        `json:"title,omitempty" bson:"title,omitempty"`
	Link     string        `json:"link,omitempty" bson:"link,omitempty"`
}
