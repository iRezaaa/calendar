package main

import (
	"gopkg.in/mgo.v2"
	"gitlab.com/irezaa/calendar/repository"
)

type Config struct {
	MongoURL string
}

type App struct {
	UserRepository    repository.UserRepository
	SessionRepository repository.SessionRepository
	EventRepository   repository.EventRepository
	NoteRepository    repository.NoteRepository
}

func (app *App) Init(config Config) {
	err, db := initDatabase(config.MongoURL)

	if err != nil {
		panic(err)
	}

	app.SessionRepository = repository.SessionRepository{}
	app.SessionRepository.DB = db

	app.UserRepository = repository.UserRepository{}
	app.UserRepository.DB = db

	app.EventRepository = repository.EventRepository{}
	app.EventRepository.DB = db

	app.NoteRepository = repository.NoteRepository{}
	app.NoteRepository.DB = db
}

func initDatabase(mongoUrl string) (error, *mgo.Database) {
	session, err := mgo.Dial(mongoUrl)

	if err != nil {
		return err, nil
	}

	session.SetMode(mgo.Monotonic, true)
	return nil, session.DB("calendar")
}
