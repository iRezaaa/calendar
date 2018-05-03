package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"gitlab.com/irezaa/calendar/model"
	"encoding/json"
)

func GetUserInfo(app *App,session *model.Session, route *Route, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := make(map[string]interface{})

	data["user_info"] = session.User

	response := HttpJsonResponse{
		ResponseStatus: ResponseStatusOk,
		Data:           data,
	}

	jsonResponse, err := json.Marshal(response)

	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}
