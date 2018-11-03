package repository

import (
	"gopkg.in/mgo.v2/bson"
	"gitlab.com/irezaa/calendar/src/model"
	"gopkg.in/mgo.v2"
)

type SessionRepository struct {
	DB *mgo.Database
}

func (r *SessionRepository) FindAll() ([]model.Session, error) {
	var sessions []model.Session
	err := r.DB.C("sessions").Find(bson.M{}).All(&sessions)
	return sessions, err
}

func (r *SessionRepository) FindByAuthToken(authToken string) (*model.Session, error) {
	var session model.Session
	err := r.DB.C("sessions").Find(bson.M{"_id": authToken}).One(&session)
	return &session, err
}

func (r *SessionRepository) FindByUser(userID string) ([]model.Session, error) {
	var sessions []model.Session
	err := r.DB.C("sessions").Find(bson.M{"user._id": userID}).All(&sessions)
	return sessions, err
}

func (r *SessionRepository) RemoveUserSessions(userID string) error {
	err := r.DB.C("sessions").Remove(bson.M{"user._id": userID})
	return err
}

func (r *SessionRepository) Insert(session *model.Session) error {
	err := r.DB.C("sessions").Insert(session)
	return err
}

func (r *SessionRepository) Delete(session model.Session) error {
	err := r.DB.C("sessions").Remove(bson.M{"_id": session.AuthToken})
	return err
}

func (r *SessionRepository) Update(session model.Session) error {
	err := r.DB.C("sessions").Update(bson.M{"_id": session.AuthToken}, &session)
	return err
}
