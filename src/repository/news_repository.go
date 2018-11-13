package repository

import (
	"gitlab.com/irezaa/calendar/src/model"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
)

type NewsRepository struct {
	DB *mgo.Database
}

func (r *NewsRepository) FindAll() ([]model.News, error) {
	var newsList []model.News
	err := r.DB.C("news").Find(bson.M{}).Sort("create_time").All(&newsList)
	return newsList, err
}

func (r *NewsRepository) FindByID(objectID bson.ObjectId) (model.News, error) {
	var news model.News

	err := r.DB.C("news").Find(bson.M{
		"_id": objectID,
	}).One(&news)

	return news, err
}

func (r *NewsRepository) Delete(news model.News) error {
	err := r.DB.C("news").Remove(bson.M{
		"_id": news.ID,
	})
	return err
}

func (r *NewsRepository) UpdateOrInsert(news model.News) (*mgo.ChangeInfo, error) {
	changeInfo, err := r.DB.C("news").Upsert(bson.M{
		"image_url":    news.ImageURL,
		"title":        news.Title,
		"content":      news.Content,
		"content_type": news.ContentType,
	}, &news)

	return changeInfo, err
}

func (r *NewsRepository) Insert(news *model.News) (*bson.ObjectId, error) {
	newID := bson.NewObjectId()
	news.ID = newID
	err := r.DB.C("news").Insert(news)

	if err != nil {
		return nil, err
	}

	return &newID, err
}

func (r *NewsRepository) Update(news *model.News) error {
	err := r.DB.C("news").Update(bson.M{"_id": news.ID}, news)
	return err
}
