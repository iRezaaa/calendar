package main

import (
	"gopkg.in/mgo.v2"
	"time"
	"fmt"
	"gitlab.com/irezaa/calendar/src/repository"
)

type Config struct {
	DBConfig DatabaseConfig
}

type DatabaseConfig struct {
	URL string
	DBName string
	UserName string
	Password string
	Timeout time.Duration
}

type App struct {
	UserRepository      repository.UserRepository
	SessionRepository   repository.SessionRepository
	EventRepository     repository.EventRepository
	NoteRepository      repository.NoteRepository
	BannerRepository    repository.BannerRepository
	IndicatorRepository repository.IndicatorRepository
}

func (app *App) Init(config Config) {
	err, db := initDatabase(config.DBConfig)

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

	app.BannerRepository = repository.BannerRepository{}
	app.BannerRepository.DB = db

	app.IndicatorRepository = repository.IndicatorRepository{}
	app.IndicatorRepository.DB = db

}

func initDatabase(dbConfig DatabaseConfig) (error, *mgo.Database) {
	info := &mgo.DialInfo{
		Addrs:    []string{dbConfig.URL},
		Timeout:  dbConfig.Timeout,
		Database: dbConfig.DBName,
		Username: dbConfig.UserName,
		Password: dbConfig.Password,
	}

	fmt.Printf("Connecting to %s" , info.Addrs[0])

	session, err := mgo.DialWithInfo(info)

	if err != nil {
		return err, nil
	}
	print("calendar_app.go::initDatabase -> Database Connected")
	session.SetMode(mgo.Monotonic, true)
	return nil, session.DB(dbConfig.DBName)
}
