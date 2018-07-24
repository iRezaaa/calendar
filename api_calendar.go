package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"gitlab.com/irezaa/calendar/model"
	"encoding/json"
	"strconv"
	"time"
	"github.com/yaa110/go-persian-calendar/ptime"
)

func GetEvents(app *App, session *model.Session, route *Route, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := make(map[string]interface{})
	var responseStatus = ResponseStatusOk

	q := r.URL.Query()
	keys := make(map[string]string)
	keys["year"] = q.Get("year")
	keys["month"] = q.Get("month")

	emptyKeys := emptyValidator(keys)

	if len(emptyKeys) != 0 {
		responseStatus = ResponseStatusError
		data["err_code"] = 2
		data["items"] = emptyKeys
	} else {
		var castError error
		var monthInt int
		var yearInt int

		monthInt, castError = strconv.Atoi(keys["month"])
		if castError != nil {
			responseStatus = ResponseStatusError
			data["err_code"] = 0
			data["err_text"] = "invalid month"
		} else {
			yearInt, castError = strconv.Atoi(keys["year"])

			if castError != nil {
				responseStatus = ResponseStatusError
				data["err_code"] = 0
				data["err_text"] = "invalid year"
			} else {
				generalEvents, getGeneralEventsError := app.EventRepository.FindByJalaliMonth(monthInt, yearInt, model.EventTypeGeneral, "")

				if getGeneralEventsError != nil {
					responseStatus = ResponseStatusError
					data["err_code"] = 0
					data["err_text"] = "general events error"
				} else {

					personalEvents, getPersonalEventsError := app.EventRepository.FindByJalaliMonth(monthInt, yearInt, model.EventTypePersonal, session.User.ID)

					if getPersonalEventsError != nil {
						responseStatus = ResponseStatusError
						data["err_code"] = 0
						data["err_text"] = getPersonalEventsError
					} else {
						panelEvents, getPanelEventsError := app.EventRepository.FindByJalaliMonth(monthInt, yearInt, model.EventTypePanel, "")

						if getPanelEventsError != nil {
							responseStatus = ResponseStatusError
							data["err_code"] = 0
							data["err_text"] = "panel events error"
						} else {
							responseStatus = ResponseStatusOk
							data["general_events"] = generalEvents
							data["personal_events"] = personalEvents
							data["panel_events"] = panelEvents
						}
					}
				}
			}
		}
	}

	response := HttpJsonResponse{
		ResponseStatus: responseStatus,
		Data:           data,
	}

	jsonResponse, err := json.Marshal(response)

	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}

func AddPersonalEvent(app *App, session *model.Session, route *Route, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := make(map[string]interface{})
	var responseStatus = ResponseStatusOk

	r.ParseMultipartForm(0)

	keys := make(map[string]string)
	keys["timeStamp"] = r.FormValue("timeStamp")
	keys["event_title"] = r.FormValue("event_title")

	emptyKeys := emptyValidator(keys)

	if len(emptyKeys) != 0 {
		responseStatus = ResponseStatusError
		data["err_code"] = 2
		data["items"] = emptyKeys
	} else {
		var castError error
		var timeStamp int64

		timeStamp, castError = strconv.ParseInt(keys["timeStamp"], 10, 64)

		gregorianTime := time.Unix(timeStamp, 0)
		jalaliTime := ptime.New(gregorianTime)

		if castError != nil {
			responseStatus = ResponseStatusError
			data["err_code"] = 0
			data["err_text"] = "invalid parameters"
		} else {
			personalEvent := model.Event{
				EventTitle:  keys["event_title"],
				EventType:   model.EventTypePersonal,
				Day:         gregorianTime.Day(),
				Month:       int(gregorianTime.Month()),
				Year:        gregorianTime.Year(),
				JalaliDay:   jalaliTime.Day(),
				JalaliMonth: int(jalaliTime.Month()),
				JalaliYear:  jalaliTime.Year(),
				UserID:      session.User.ID,
				IsHoliday:   false,
			}

			err := app.EventRepository.UpdateOrInsert(personalEvent)

			if err != nil {
				responseStatus = ResponseStatusError
				data["err_code"] = 0
			} else {
				responseStatus = ResponseStatusOk
			}
		}
	}

	response := HttpJsonResponse{
		ResponseStatus: responseStatus,
		Data:           data,
	}

	jsonResponse, err := json.Marshal(response)

	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}

func DeletePersonalEvent(app *App, session *model.Session, route *Route, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := make(map[string]interface{})
	var responseStatus = ResponseStatusOk

	r.ParseMultipartForm(0)

	keys := make(map[string]string)
	keys["event_id"] = r.FormValue("event_id")

	emptyKeys := emptyValidator(keys)

	if len(emptyKeys) != 0 {
		responseStatus = ResponseStatusError
		data["err_code"] = 2
		data["items"] = emptyKeys
	} else {
		event, err := app.EventRepository.FindByID(keys["event_id"], model.EventTypePersonal)

		if err != nil {
			responseStatus = ResponseStatusError
			data["err_code"] = 0
			data["err_text"] = "event not found"
		} else {
			if event.UserID == session.User.ID {
				err = app.EventRepository.Delete(event)

				if err != nil {
					responseStatus = ResponseStatusError
					data["err_code"] = 0
					data["err_text"] = "cant delete personal event!"
				} else {
					responseStatus = ResponseStatusOk
				}
			} else {
				responseStatus = ResponseStatusError
				data["err_code"] = 1
				data["err_text"] = "you are not event owner :D"
			}

		}
	}

	response := HttpJsonResponse{
		ResponseStatus: responseStatus,
		Data:           data,
	}

	jsonResponse, err := json.Marshal(response)

	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}

func AddPanelEvent(app *App, session *model.Session, route *Route, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := make(map[string]interface{})
	var responseStatus = ResponseStatusOk

	r.ParseMultipartForm(0)

	keys := make(map[string]string)
	keys["timeStamp"] = r.FormValue("timeStamp")
	keys["event_title"] = r.FormValue("event_title")

	emptyKeys := emptyValidator(keys)

	if len(emptyKeys) != 0 {
		responseStatus = ResponseStatusError
		data["err_code"] = 2
		data["items"] = emptyKeys
	} else {
		var castError error
		var timeStamp int64

		timeStamp, castError = strconv.ParseInt(keys["timeStamp"], 10, 64)

		gregorianTime := time.Unix(timeStamp, 0)
		jalaliTime := ptime.New(gregorianTime)

		if castError != nil {
			responseStatus = ResponseStatusError
			data["err_code"] = 0
			data["err_text"] = "invalid parameters"
		} else {
			panelEvent := model.Event{
				EventTitle:  keys["event_title"],
				EventType:   model.EventTypePanel,
				Day:         gregorianTime.Day(),
				Month:       int(gregorianTime.Month()),
				Year:        gregorianTime.Year(),
				JalaliDay:   jalaliTime.Day(),
				JalaliMonth: int(jalaliTime.Month()),
				JalaliYear:  jalaliTime.Year(),
				UserID:      session.User.ID,
				IsHoliday:   false,
			}

			err := app.EventRepository.UpdateOrInsert(panelEvent)

			if err != nil {
				responseStatus = ResponseStatusError
				data["err_code"] = 0
			} else {
				responseStatus = ResponseStatusOk
			}
		}
	}

	response := HttpJsonResponse{
		ResponseStatus: responseStatus,
		Data:           data,
	}

	jsonResponse, err := json.Marshal(response)

	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}

func DeletePanelEvent(app *App, session *model.Session, route *Route, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := make(map[string]interface{})
	var responseStatus = ResponseStatusOk

	r.ParseMultipartForm(0)

	keys := make(map[string]string)
	keys["event_id"] = r.FormValue("event_id")

	emptyKeys := emptyValidator(keys)

	if len(emptyKeys) != 0 {
		responseStatus = ResponseStatusError
		data["err_code"] = 2
		data["items"] = emptyKeys
	} else {
		event, err := app.EventRepository.FindByID(keys["event_id"], model.EventTypePanel)

		if err != nil {
			responseStatus = ResponseStatusError
			data["err_code"] = 0
			data["err_text"] = "event not found"
		} else {
			err = app.EventRepository.Delete(event)

			if err != nil {
				responseStatus = ResponseStatusError
				data["err_code"] = 0
				data["err_text"] = "cant delete panel event!"
			} else {
				responseStatus = ResponseStatusOk
			}
		}
	}

	response := HttpJsonResponse{
		ResponseStatus: responseStatus,
		Data:           data,
	}

	jsonResponse, err := json.Marshal(response)

	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}
