package repository

import (
	"gitlab.com/irezaa/calendar/src/model"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
)

type BannerRepository struct {
	DB *mgo.Database
}

func (r *BannerRepository) FindAll() ([]model.Banner, error) {
	var banners []model.Banner
	err := r.DB.C("banners").Find(bson.M{}).All(&banners)
	return banners, err
}

func (r *BannerRepository) FindByID(objectID bson.ObjectId) (model.Banner, error) {
	var banner model.Banner

	err := r.DB.C("banners").Find(bson.M{
		"_id": objectID,
	}).One(&banner)

	return banner, err
}

func (r *BannerRepository) Delete(banner model.Banner) error {
	err := r.DB.C("banners").Remove(bson.M{
		"_id": banner.ID,
	})
	return err
}

func (r *BannerRepository) UpdateOrInsert(banner model.Banner) (*mgo.ChangeInfo, error) {
	changeInfo, err := r.DB.C("banners").Upsert(bson.M{
		"image_url": banner.ImageURL,
		"title":     banner.Title,
		"link":      banner.Link,
	}, &banner)

	return changeInfo, err
}
