package model

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Day struct {
	ID             bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	YearNumber     int           `json:"year_number" bson:"year_number"`
	MonthNumber    int           `json:"month_number" bson:"month_number"`
	StartTimeStamp time.Time     `json:"start_timestamp" bson:"start_timestamp"`
	EndTimeStamp   time.Time     `json:"end_timestamp" bson:"end_timestamp"`
}