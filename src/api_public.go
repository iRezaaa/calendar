package src

import (
	"gitlab.com/irezaa/calendar/src/model"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"encoding/json"
	"time"
	"fmt"
	"strconv"
	"github.com/yaa110/go-persian-calendar/ptime"
)

func Today(app *App, session *model.Session, route *Route, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := make(map[string]interface{})
	var responseStatus = ResponseStatusOk

	now := time.Now()

	fmt.Printf("Month : %d - Year : %d - Day : %d", int(now.Month()), now.Year(), now.Day())

	var isHoliday = false

	generalEvents, err := app.EventRepository.FindByGregorianDay(int(now.Month()), now.Year(), now.Day(), model.EventTypeGeneral, "")
	panelEvents, err := app.EventRepository.FindByGregorianDay(int(now.Month()), now.Year(), now.Day(), model.EventTypePanel, "")

	var personalEvents []model.Event
	var todayNotes []model.Note

	if session != nil && session.User != nil && session.User.ID != "" {
		personalEvents, err = app.EventRepository.FindByGregorianDay(int(now.Month()), now.Year(), now.Day(), model.EventTypePersonal, session.User.ID)
		todayNotes, err = app.NoteRepository.FindLastByUserID(session.User.ID, now.Year(), int(now.Month()), now.Day())
	}


	if generalEvents != nil && len(generalEvents) > 0 {
		for _, event := range generalEvents {
			if event.IsHoliday {
				isHoliday = true
				break
			}
		}
	} else {
		isHoliday = false
	}

	if err == nil {
		data["general_events"] = generalEvents

		if session != nil && session.User != nil && session.User.ID != "" {
			data["personal_events"] = personalEvents

			var todayNote *model.Note

			if todayNotes != nil && len(todayNotes) > 0 {
				todayNote = &todayNotes[0]
			}else{
				todayNote = nil
			}

			data["today_note"] = todayNote
		}

		data["panel_events"] = panelEvents
		data["isHoliday"] = isHoliday
	} else {
		print(err.Error())
		responseStatus = ResponseStatusError
		data["err_code"] = 0
		data["err_text"] = "error while getting events"
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

func GetDay(app *App, session *model.Session, route *Route, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := make(map[string]interface{})
	var responseStatus = ResponseStatusOk

	q := r.URL.Query()
	keys := make(map[string]string)
	keys["year"] = q.Get("year")
	keys["month"] = q.Get("month")
	keys["day"] = q.Get("day")
	emptyKeys := emptyValidator(keys)

	if emptyKeys != nil && len(emptyKeys) != 0 {
		responseStatus = ResponseStatusError
		data["err_code"] = 2
		data["items"] = emptyKeys
	} else {
		var castError error
		var yearInt int
		var monthInt int
		var dayInt int

		yearInt, castError = strconv.Atoi(keys["year"])

		if castError != nil {
			responseStatus = ResponseStatusError
			data["err_code"] = 0
			data["err_text"] = "invalid year"
		} else {

			monthInt, castError = strconv.Atoi(keys["month"])

			if castError != nil {
				responseStatus = ResponseStatusError
				data["err_code"] = 0
				data["err_text"] = "invalid month"
			} else {
				dayInt, castError = strconv.Atoi(keys["day"])

				if castError != nil {
					responseStatus = ResponseStatusError
					data["err_code"] = 0
					data["err_text"] = "invalid day"
				} else {
					var pt = ptime.Date(yearInt, ptime.Month(monthInt), dayInt, 0, 0, 0, 0, ptime.Iran())

					jalaliTimeStamp := pt.Time().Unix()
					gregorianTime := time.Unix(jalaliTimeStamp, 0).In(ptime.Iran())

					fmt.Printf("Month : %d - Year : %d - Day : %d", int(gregorianTime.Month()), gregorianTime.Year(), gregorianTime.Day())

					var isHoliday = false

					generalEvents, err := app.EventRepository.FindByGregorianDay(int(gregorianTime.Month()), gregorianTime.Year(), gregorianTime.Day(), model.EventTypeGeneral, "")
					panelEvents, err := app.EventRepository.FindByGregorianDay(int(gregorianTime.Month()), gregorianTime.Year(), gregorianTime.Day(), model.EventTypePanel, "")

					var personalEvents []model.Event
					var todayNotes []model.Note

					if session != nil && session.User != nil && session.User.ID != "" {
						personalEvents, err = app.EventRepository.FindByGregorianDay(int(gregorianTime.Month()), gregorianTime.Year(), gregorianTime.Day(), model.EventTypePersonal, session.User.ID)
						todayNotes, err = app.NoteRepository.FindLastByUserID(session.User.ID, gregorianTime.Year(), int(gregorianTime.Month()), gregorianTime.Day())
					}

					if generalEvents != nil && len(generalEvents) > 0 {
						for _, event := range generalEvents {
							if event.IsHoliday {
								isHoliday = true
								break
							}
						}
					} else {
						isHoliday = false
					}

					if err == nil {
						data["general_events"] = generalEvents

						if session != nil && session.User != nil && session.User.ID != "" {
							data["personal_events"] = personalEvents

							var todayNote *model.Note

							if todayNotes != nil && len(todayNotes) > 0 {
								todayNote = &todayNotes[0]
							}else{
								todayNote = nil
							}

							data["today_note"] = todayNote
						}

						data["panel_events"] = panelEvents
						data["isHoliday"] = isHoliday
					} else {
						print(err.Error())
						responseStatus = ResponseStatusError
						data["err_code"] = 0
						data["err_text"] = "error while getting events"
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
