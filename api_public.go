package main

import (
	"gitlab.com/irezaa/calendar/model"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"encoding/json"
	"time"
	"fmt"
)

func Today(app *App, session *model.Session, route *Route, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := make(map[string]interface{})
	var responseStatus = ResponseStatusOk

	now := time.Now()

	fmt.Printf("Month : %d - Year : %d - Day : %d", int(now.Month()), now.Year(), now.Day())

	var isHoliday = false

	generalEvents, err := app.EventRepository.FindByGregorianDay(int(now.Month()), now.Year(), now.Day(), model.EventTypeGeneral, session.User.ID)
	personalEvents, err := app.EventRepository.FindByGregorianDay(int(now.Month()), now.Year(), now.Day(), model.EventTypePersonal, session.User.ID)
	panelEvents, err := app.EventRepository.FindByGregorianDay(int(now.Month()), now.Year(), now.Day(), model.EventTypePanel, session.User.ID)
	todayNotes, err := app.NoteRepository.FindLastByUserID(session.User.ID, now.Year(), int(now.Month()), now.Day())

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
		data["personal_events"] = personalEvents
		data["panel_events"] = panelEvents

		if todayNotes != nil && len(todayNotes) > 0 {
			data["today_note"] = todayNotes[0]
		}

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
