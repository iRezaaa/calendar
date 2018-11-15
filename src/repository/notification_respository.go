package repository

import (
"gopkg.in/mgo.v2"
"gopkg.in/mgo.v2/bson"
"gitlab.com/irezaa/calendar/src/model"
)

type NotificationRepository struct {
	DB *mgo.Database
}

func (r *NotificationRepository) FindAll() ([]model.Notification, error) {
	var notifications []model.Notification
	err := r.DB.C("notifications").Find(bson.M{}).All(&notifications)
	return notifications, err
}

func (r *NotificationRepository) FindByID(objectID bson.ObjectId) (model.Notification, error) {
	var notification model.Notification

	err := r.DB.C("notifications").Find(bson.M{
		"_id": objectID,
	}).One(&notification)

	return notification, err
}

func (r *NotificationRepository) Delete(note model.Notification) error {
	err := r.DB.C("notifications").Remove(bson.M{
		"_id": note.ID,
	})
	return err
}

func (r *NotificationRepository) Insert(notification *model.Notification) (*bson.ObjectId, error) {
	newID := bson.NewObjectId()
	notification.ID = newID
	err := r.DB.C("notifications").Insert(notification)

	if err != nil {
		return nil, err
	}

	return &newID, err
}
