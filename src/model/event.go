package model

import "gopkg.in/mgo.v2/bson"

type EventType int

const (
	EventTypeGeneral  EventType = 0x0001
	EventTypePanel    EventType = 0x0002
	EventTypePersonal EventType = 0x0003
)

type Event struct {
	ID          bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	EventTitle  string        `json:"eventTitle,omitempty" bson:"eventTitle,omitempty"`
	EventType   EventType     `json:"eventType,omitempty" bson:"eventType,omitempty"`
	JalaliYear  int           `json:"jalaliYear,omitempty" bson:"jalaliYear,omitempty"`
	JalaliMonth int           `json:"jalaliMonth,omitempty" bson:"jalaliMonth,omitempty"`
	JalaliDay   int           `json:"jalaliDay,omitempty" bson:"jalaliDay,omitempty"`
	Day         int           `json:"gregorianDay,omitempty" bson:"gregorianDay,omitempty"`
	Month       int           `json:"gregorianMonth,omitempty" bson:"gregorianMonth,omitempty"`
	Year        int           `json:"gregorianYear,omitempty" bson:"gregorianYear,omitempty"`
	UserID      string        `json:"-" bson:"user_id,omitempty"`
	IsHoliday   bool          `json:"isHoliday" bson:"isHoliday"`
}
