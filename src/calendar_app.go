package main

import (
	"gopkg.in/mgo.v2"
	"time"
	"fmt"
	"gitlab.com/irezaa/calendar/src/repository"
	"github.com/NaySoftware/go-fcm"
)

type Config struct {
	DBConfig  DatabaseConfig
	FCMConfig FCMConfig
}

type FCMConfig struct {
	ServerKey    string
	DefaultTopic string
}

type DatabaseConfig struct {
	URL      string
	DBName   string
	UserName string
	Password string
	Timeout  time.Duration
}

type App struct {
	UserRepository         repository.UserRepository
	SessionRepository      repository.SessionRepository
	EventRepository        repository.EventRepository
	NoteRepository         repository.NoteRepository
	BannerRepository       repository.BannerRepository
	IndicatorRepository    repository.IndicatorRepository
	NewsRepository         repository.NewsRepository
	NotificationRepository repository.NotificationRepository
	PushService            PushService
}

func (app *App) Init(config Config) {
	err, db := initDatabase(config.DBConfig)

	if err != nil {
		panic(err)
	}

	fcmClient := initFCM(config.FCMConfig)

	if fcmClient != nil {
		app.PushService = PushService{}
		app.PushService.FCMClient = fcmClient
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

	app.NewsRepository = repository.NewsRepository{}
	app.NewsRepository.DB = db

	app.NotificationRepository = repository.NotificationRepository{}
	app.NotificationRepository.DB = db
}

func initFCM(fcmConfig FCMConfig) *fcm.FcmClient {
	return fcm.NewFcmClient(fcmConfig.ServerKey)
}

func initDatabase(dbConfig DatabaseConfig) (error, *mgo.Database) {
	info := &mgo.DialInfo{
		Addrs:    []string{dbConfig.URL},
		Timeout:  dbConfig.Timeout,
		Database: dbConfig.DBName,
		Username: dbConfig.UserName,
		Password: dbConfig.Password,
	}

	fmt.Printf("Connecting to %s", info.Addrs[0])

	session, err := mgo.DialWithInfo(info)

	if err != nil {
		return err, nil
	}
	print("calendar_app.go::initDatabase -> Database Connected")
	session.SetMode(mgo.Monotonic, true)
	return nil, session.DB(dbConfig.DBName)
}
