package model

type UserType int

const (
	UserTypeNormal UserType = iota + 1
	UserTypeAdmin
)

type User struct {
	ID       string    `json:"id,omitempty" bson:"_id,omitempty"`
	FName    string    `json:"f_name" bson:"f_name"`
	LName    string    `json:"l_name" bson:"l_name"`
	Pass     string    `json:"-" bson:"password"`
	Type     UserType  `json:"type" bson:"type"`
}
