package repository

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gitlab.com/irezaa/calendar/model"
)

type EventRepository struct {
	DB *mgo.Database
}

func (r *EventRepository) FindByID(eventID string, eventType model.EventType) (model.Event, error) {
	var event model.Event
	var err error

	switch eventType {
	case model.EventTypeGeneral:
		err = r.DB.C("events_general").FindId(bson.ObjectIdHex(eventID)).One(&event)
		break
	case model.EventTypePersonal:
		err = r.DB.C("events_personal").FindId(bson.ObjectIdHex(eventID)).One(&event)
		break
	case model.EventTypePanel:
		err = r.DB.C("events_panel").FindId(bson.ObjectIdHex(eventID)).One(&event)
		break
	}

	return event, err
}

func (r *EventRepository) FindByJalaliMonth(jalaliMonth int, jalaliYear int, eventType model.EventType, userID string) ([]model.Event, error) {
	var events []model.Event
	var err error

	switch eventType {
	case model.EventTypeGeneral:
		err = r.DB.C("events_general").Find(bson.M{
			"jalaliMonth": jalaliMonth,
			"jalaliYear":  jalaliYear,
		}).All(&events)
		break
	case model.EventTypePersonal:
		err = r.DB.C("events_personal").Find(bson.M{
			"jalaliMonth": jalaliMonth,
			"jalaliYear":  jalaliYear,
			"user_id":     userID,
		}).All(&events)
		break
	case model.EventTypePanel:
		err = r.DB.C("events_panel").Find(bson.M{
			"jalaliMonth": jalaliMonth,
			"jalaliYear":  jalaliYear,
		}).All(&events)
		break
	}

	return events, err
}

func (r *EventRepository) FindByGregorianMonth(gregorianMonth int, gregorianYear int, eventType model.EventType, userID string) ([]model.Event, error) {
	var events []model.Event
	var err error

	switch eventType {
	case model.EventTypeGeneral:
		err = r.DB.C("events_general").Find(bson.M{
			"gregorianMonth": gregorianMonth,
			"gregorianYear":  gregorianYear,
		}).All(&events)
		break
	case model.EventTypePersonal:
		err = r.DB.C("events_personal").Find(bson.M{
			"gregorianMonth": gregorianMonth,
			"gregorianYear":  gregorianYear,
			"user_id":        userID,
		}).All(&events)
		break
	case model.EventTypePanel:
		err = r.DB.C("events_panel").Find(bson.M{
			"gregorianMonth": gregorianMonth,
			"gregorianYear":  gregorianYear,
		}).All(&events)
		break
	}

	return events, err
}

func (r *EventRepository) FindByGregorianDay(gregorianMonth int, gregorianYear int, gregorianDay int, eventType model.EventType, userID string) ([]model.Event, error) {
	var events []model.Event
	var err error

	switch eventType {
	case model.EventTypeGeneral:
		err = r.DB.C("events_general").Find(bson.M{
			"gregorianMonth": gregorianMonth,
			"gregorianYear":  gregorianYear,
			"gregorianDay":   gregorianDay,
		}).All(&events)
		break
	case model.EventTypePersonal:
		err = r.DB.C("events_personal").Find(bson.M{
			"gregorianMonth": gregorianMonth,
			"gregorianYear":  gregorianYear,
			"gregorianDay":   gregorianDay,
			"user_id":        userID,
		}).All(&events)
		break
	case model.EventTypePanel:
		err = r.DB.C("events_panel").Find(bson.M{
			"gregorianMonth": gregorianMonth,
			"gregorianYear":  gregorianYear,
			"gregorianDay":   gregorianDay,
		}).All(&events)
		break
	}

	return events, err
}

func (r *EventRepository) FindTodayPersonalEvents(gregorianMonth int, gregorianYear int, gregorianDay int) ([]model.Event, error) {
	var personalEvents []model.Event
	var err error

	err = r.DB.C("events_personal").Find(bson.M{
		"gregorianMonth": gregorianMonth,
		"gregorianYear":  gregorianYear,
		"gregorianDay":   gregorianDay,
	}).All(&personalEvents)

	return personalEvents, err
}

func (r *EventRepository) Delete(event model.Event) error {
	var err error

	switch event.EventType {
	case model.EventTypeGeneral:
		err = r.DB.C("events_general").Remove(bson.M{
			"_id": event.ID,
		})
		break
	case model.EventTypePersonal:
		err = r.DB.C("events_personal").Remove(bson.M{
			"_id": event.ID,
		})
		break
	case model.EventTypePanel:
		err = r.DB.C("events_panel").Remove(bson.M{
			"_id": event.ID,
		})
		break
	}

	return err
}

func (r *EventRepository) UpdateOrInsert(event model.Event) error {
	var err error

	switch event.EventType {
	case model.EventTypeGeneral:
		_, err = r.DB.C("events_general").Upsert(bson.M{
			"eventType":      event.EventType,
			"gregorianDay":   event.Day,
			"gregorianMonth": event.Month,
			"gregorianYear":  event.Year,
			"eventTitle":     event.EventTitle,
		}, &event)

		break
	case model.EventTypePersonal:
		_, err = r.DB.C("events_personal").Upsert(bson.M{
			"eventType":      event.EventType,
			"gregorianDay":   event.Day,
			"gregorianMonth": event.Month,
			"gregorianYear":  event.Year,
			"eventTitle":     event.EventTitle,
		}, &event)
		break
	case model.EventTypePanel:
		_, err = r.DB.C("events_panel").Upsert(bson.M{
			"eventType":      event.EventType,
			"gregorianDay":   event.Day,
			"gregorianMonth": event.Month,
			"gregorianYear":  event.Year,
			"eventTitle":     event.EventTitle,
		}, &event)
		break
	}

	return err
}
