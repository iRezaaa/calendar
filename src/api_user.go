package src

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"gitlab.com/irezaa/calendar/src/model"
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
	data := make(map[string]interface{})
	var responseStatus = ResponseStatusOk

	r.ParseMultipartForm(0)

	keys := make(map[string]string)
	keys["username"] = r.FormValue("username")
	emptyKeys := emptyValidator(keys)

	if len(emptyKeys) != 0 {
		responseStatus = ResponseStatusError
		data["err_code"] = 3
		data["items"] = emptyKeys
	} else {
		user, err := app.UserRepository.FindByID(keys["username"])

		if err == nil && user != nil {
			user.Type = model.UserTypeAdmin
			err = app.UserRepository.Update(user)

			if err != nil {
				responseStatus = ResponseStatusError
				data["err_code"] = 0
				data["err_text"] = "cannot update user model to admin"
				data["server_error"] = err.Error()
			} else {
				err := app.SessionRepository.RemoveUserSessions(user.ID)

				if err != nil {
					responseStatus = ResponseStatusError
					data["user_info"] = user
					data["err_code"] = 0
					data["err_text"] = "cannot delete sessions"
				} else {
					responseStatus = ResponseStatusOk
					data["user_info"] = user
				}
			}
		} else {
			responseStatus = ResponseStatusError
			data["err_code"] = 1
			data["err_text"] = "cannot find user"
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
