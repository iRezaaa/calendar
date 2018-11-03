package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"time"
	"encoding/json"
	"gitlab.com/irezaa/calendar/src/model"
		"golang.org/x/crypto/bcrypt"
)

func Login(app *App, session *model.Session, route *Route, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := make(map[string]interface{})
	var responseStatus = ResponseStatusOk

	r.ParseMultipartForm(0)
	keys := make(map[string]string)
	keys["username"] = r.FormValue("username")
	keys["password"] = r.FormValue("password")
	keys["fcm_token"] = r.FormValue("fcm_token")

	emptyKeys := emptyValidator(keys)

	if len(emptyKeys) == 0 {
		user, err := app.UserRepository.FindByID(keys["username"])

		if err == nil {
			passwordFromDatabase := []byte(user.Pass)
			passFromRequest := []byte(keys["password"])

			err = bcrypt.CompareHashAndPassword(passwordFromDatabase, passFromRequest)

			if err == nil {
				authToken := RandString(15)

				newSession := model.Session{
					AuthToken: authToken,
					FcmToken:  keys["fcm_token"],
					OpenTime:  time.Now(),
					User:      user,
				}

				err := app.SessionRepository.Insert(&newSession)

				if err == nil {
					responseStatus = ResponseStatusOk
					data["auth_token"] = authToken
				}else{
					responseStatus = ResponseStatusError
					data["err_code"] = 0
					data["err_text"] = "can not create session"
				}
			}else{
				responseStatus = ResponseStatusError
				data["err_code"] = 2
				data["err_text"] = "password is wrong!"
			}
		}else{
			responseStatus = ResponseStatusError
			data["err_code"] = 1
			data["err_text"] = "this username and password not registered!"
		}
	} else {
		responseStatus = ResponseStatusError
		data["err_code"] = 3
		data["items"] = emptyKeys
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

func Register(app *App, session *model.Session, route *Route, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := make(map[string]interface{})
	var responseStatus = ResponseStatusOk

	r.ParseMultipartForm(0)

	keys := make(map[string]string)

	keys["username"] = r.FormValue("username")
	keys["password"] = r.FormValue("password")
	keys["first_name"] = r.FormValue("first_name")
	keys["last_name"] = r.FormValue("last_name")
	keys["fcm_token"] = r.FormValue("fcm_token")

	emptyKeys := emptyValidator(keys)

	if len(emptyKeys) != 0 {
		responseStatus = ResponseStatusError
		data["err_code"] = 3
		data["items"] = emptyKeys
		data["err_text"] = "missing required keys"
	}else{
		authToken := RandString(15)
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(keys["password"]), bcrypt.DefaultCost)

		if err != nil {
			responseStatus = ResponseStatusError
			data["err_code"] = 0
			data["err_text"] = "password cannot be hashed!"
		}else {
			user := new(model.User)
			user.ID = keys["username"]
			user.FName = keys["first_name"]
			user.LName = keys["last_name"]
			user.Pass = string(hashedPassword)
			user.Type = model.UserTypeNormal

			err := app.UserRepository.Insert(user)

			if err != nil {
				responseStatus = ResponseStatusError
				data["err_code"] = 0
				data["err_text"] = "cannot insert user object to database!"
			}else{
				session := model.Session{
					AuthToken: authToken,
					FcmToken:  keys["fcm_token"],
					OpenTime:  time.Now(),
					User:      user,
				}

				err := app.SessionRepository.Insert(&session)

				if err != nil {
					responseStatus = ResponseStatusError
					data["err_code"] = 0
					data["err_text"] = "cannot insert session object to database!"
				}else{
					responseStatus = ResponseStatusOk
					data["auth_token"] = authToken
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
