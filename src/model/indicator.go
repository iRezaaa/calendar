package model

type Indicator struct {
	ID      string `json:"id,omitempty" bson:"id,omitempty"`
	Title   string `json:"title,omitempty" bson:"title,omitempty"`
	Value   string `json:"value,omitempty" bson:"value,omitempty"`
	Percent string `json:"percent,omitempty" bson:"percent,omitempty"`
}
