package model

import "time"

type Session struct {
	AuthToken string    `json:"auth_token,omitempty" bson:"_id,omitempty"`
	FcmToken  string    `json:"fcm_token" bson:"fcm_token"`
	User      *User      `json:"user,omitempty" bson:"user,omitempty"`
	OpenTime  time.Time `json:"open_time" bson:"open_time"`
}
