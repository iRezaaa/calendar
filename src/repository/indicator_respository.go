package repository

import (
	"gitlab.com/irezaa/calendar/src/model"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
)

type IndicatorRepository struct {
	DB *mgo.Database
}

func (r *IndicatorRepository) FindAll() ([]model.Indicator, error) {
	var indicators []model.Indicator
	err := r.DB.C("indicators").Find(bson.M{}).All(&indicators)
	return indicators, err
}

func (r *IndicatorRepository) FindByID(id string) (model.Indicator, error) {
	var indicator model.Indicator

	err := r.DB.C("indicators").Find(bson.M{
		"id": id,
	}).One(&indicator)

	return indicator, err
}

func (r *IndicatorRepository) Delete(indicator model.Indicator) error {
	err := r.DB.C("indicators").Remove(bson.M{
		"id": indicator.ID,
	})
	return err
}

func (r *IndicatorRepository) UpdateOrInsert(indicator model.Indicator) (*mgo.ChangeInfo, error) {
	changeInfo, err := r.DB.C("indicators").Upsert(bson.M{
		"id":    indicator.ID,
	}, &indicator)

	return changeInfo, err
}
