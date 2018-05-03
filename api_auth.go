package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"time"
	"encoding/json"
	"gitlab.com/irezaa/calendar/model"
	"fmt"
)

func Login(app *App, session *model.Session, route *Route, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseMultipartForm(0)

	keys := make(map[string]string)
	keys["username"] = r.FormValue("username")
	keys["password"] = r.FormValue("password")
	keys["fcm_token"] = r.FormValue("fcm_token")

	emptyKeys := emptyValidator(keys)

	if len(emptyKeys) == 0 {
		user, err := app.UserRepository.FindByID(keys["username"])

		if err == nil {
			if user.Pass == keys["password"] {
				authToken := RandString(15)

				newSession := model.Session{
					AuthToken: authToken,
					FcmToken:  keys["fcm_token"],
					OpenTime:  time.Now(),
					User:      user,
				}

				err := app.SessionRepository.Insert(&newSession)

				if err == nil {
					data := make(map[string]interface{})

					data["auth_token"] = authToken

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

func Register(app *App, session *model.Session, route *Route, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseMultipartForm(0)

	keys := make(map[string]string)

	keys["username"] = r.FormValue("username")
	keys["password"] = r.FormValue("password")
	keys["first_name"] = r.FormValue("first_name")
	keys["last_name"] = r.FormValue("last_name")
	keys["fcm_token"] = r.FormValue("fcm_token")

	emptyKeys := emptyValidator(keys)

	if len(emptyKeys) == 0 {

		authToken := RandString(15)

		user := new(model.User)
		user.ID = keys["username"]
		user.FName = keys["first_name"]
		user.LName = keys["last_name"]
		user.Pass = keys["password"]
		user.Type = model.UserTypeNormal

		err := app.UserRepository.Insert(user)

		if err == nil {
			session := model.Session{
				AuthToken: authToken,
				FcmToken:  keys["fcm_token"],
				OpenTime:  time.Now(),
				User:      *user,
			}

			err := app.SessionRepository.Insert(&session)

			if err == nil {
				data := make(map[string]interface{})

				data["auth_token"] = authToken

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
				fmt.Print(err)
			}
		}else{
			fmt.Print(err)
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
