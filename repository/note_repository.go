package repository

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gitlab.com/irezaa/calendar/model"
)

type NoteRepository struct {
	DB *mgo.Database
}

func (r *NoteRepository) FindAll() ([]model.Note, error) {
	var notes []model.Note
	err := r.DB.C("notes").Find(bson.M{}).All(&notes)
	return notes, err
}

func (r *NoteRepository) FindByID(noteID string) (model.Note, error) {
	var note model.Note

	err := r.DB.C("notes").Find(bson.M{"_id": noteID}).One(&note)

	return note, err
}

func (r *NoteRepository) FindByUserID(userID string) ([]model.Note, error) {
	var notes []model.Note
	err := r.DB.C("notes").Find(bson.M{
		"user_id": userID,
	}).All(&notes)
	return notes, err
}

func (r *NoteRepository) FindLastByUserID(userID string, gregorianYear int, gregorianMonth int, gregorianDay int) ([]model.Note, error) {
	var note []model.Note
	err := r.DB.C("notes").Find(bson.M{
		"user_id": userID,
		"day":     gregorianDay,
		"month":   gregorianMonth,
		"year":    gregorianYear,
	}).Limit(1).All(&note)

	return note, err
}

func (r *NoteRepository) Delete(note model.Note) error {
	err := r.DB.C("notes").Remove(bson.M{
		"_id": note.ID,
	})
	return err
}

func (r *NoteRepository) UpdateOrInsert(note model.Note) error {
	_, err := r.DB.C("notes").Upsert(bson.M{
		"day":     note.Day,
		"month":   note.Month,
		"year":    note.Year,
		"user_id": note.UserID,
	}, &note)

	return err
}
