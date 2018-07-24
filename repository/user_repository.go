package repository

import (
	"gopkg.in/mgo.v2/bson"
	"gitlab.com/irezaa/calendar/model"
	"gopkg.in/mgo.v2"
)

type UserRepository struct {
	DB *mgo.Database
}

func (r *UserRepository) FindAll() ([]model.User, error) {
	var users []model.User
	err := r.DB.C("users").Find(bson.M{}).All(&users)
	return users, err
}

func (r *UserRepository) FindByID(userID string) (*model.User, error) {
	var user model.User

	err := r.DB.C("users").Find(bson.M{"_id" : userID}).One(&user)
	return &user, err
}

func (r *UserRepository) Insert(user *model.User) error {
	err := r.DB.C("users").Insert(user)
	return err
}

func (r *UserRepository) Delete(user model.User) error {
	err := r.DB.C("users").Remove(bson.M{"_id" : user.ID})
	return err
}

func (r *UserRepository) Update(user *model.User) error {
	err := r.DB.C("users").Update(bson.M{"_id" : user.ID},user)
	return err
}
