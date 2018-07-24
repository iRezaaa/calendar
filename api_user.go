package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"gitlab.com/irezaa/calendar/model"
	"encoding/json"
)

func GetUserInfo(app *App, session *model.Session, route *Route, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

func GrantUserToAdmin(app *App, session *model.Session, route *Route, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseMultipartForm(0)

	keys := make(map[string]string)

	keys["username"] = r.FormValue("username")

	emptyKeys := emptyValidator(keys)

	if len(emptyKeys) <= 0 {
		user, err := app.UserRepository.FindByID(keys["username"])

		if err == nil && user != nil {
			user.Type = model.UserTypeAdmin
			err = app.UserRepository.Update(user)

			if err == nil {
				app.SessionRepository.RemoveUserSessions(user.ID)

				data := make(map[string]interface{})

				data["user_info"] = user

				response := HttpJsonResponse{
					ResponseStatus: ResponseStatusOk,
					Data:           data,
				}

				jsonResponse, err := json.Marshal(response)

				if err == nil {
					w.Header().Set("Content-Type", "application/json")
					w.Write(jsonResponse)
				}
			}else{
				data := make(map[string]interface{})

				data["err_code"] = 1
				data["items"] = []string{"username"}

				response := HttpJsonResponse{
					ResponseStatus: ResponseStatusError,
					Data:           data,
				}

				jsonResponse, err := json.Marshal(response)

				if err == nil {
					w.Header().Set("Content-Type", "application/json")
					w.Write(jsonResponse)
				}
			}
		} else {
			data := make(map[string]interface{})

			data["err_code"] = 1
			data["items"] = []string{"username"}

			response := HttpJsonResponse{
				ResponseStatus: ResponseStatusError,
				Data:           data,
			}

			jsonResponse, err := json.Marshal(response)

			if err == nil {
				w.Header().Set("Content-Type", "application/json")
				w.Write(jsonResponse)
			}
		}
	} else {
		data := make(map[string]interface{})

		data["err_code"] = 3
		data["items"] = emptyKeys

		response := HttpJsonResponse{
			ResponseStatus: ResponseStatusError,
			Data:           data,
		}

		jsonResponse, err := json.Marshal(response)

		if err == nil {
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonResponse)
		}
	}

}
