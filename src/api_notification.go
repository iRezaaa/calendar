package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"encoding/json"
	"gitlab.com/irezaa/calendar/src/model"
	"time"
	"encoding/hex"
	"gopkg.in/mgo.v2/bson"
)

func GetNotificationHistory(app *App, session *model.Session, route *Route, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := make(map[string]interface{})
	var responseStatus = ResponseStatusOk

	notifications, err := app.NotificationRepository.FindAll()

	if err != nil {
		responseStatus = ResponseStatusError
		data["err_code"] = 0
		data["err_server"] = err.Error()
		data["err_text"] = "error while fetch notifications from database!"
	} else {
		responseStatus = ResponseStatusOk
		data["notifications"] = notifications
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

func PushNotification(app *App, session *model.Session, route *Route, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := make(map[string]interface{})
	var responseStatus = ResponseStatusOk

	r.ParseMultipartForm(0)

	keys := make(map[string]string)
	keys["title"] = r.FormValue("title")
	keys["body"] = r.FormValue("body")

	emptyKeys := emptyValidator(keys)

	if len(emptyKeys) != 0 {
		responseStatus = ResponseStatusError
		data["err_code"] = 2
		data["items"] = emptyKeys
	} else {
		status, err := app.PushService.PushToAll("panel", map[string]string{
			"title": keys["title"],
			"body":  keys["body"],
		})

		if err != nil {
			responseStatus = ResponseStatusError
			data["err_code"] = 2
			data["err_text"] = "push error"
			data["err_server"] = err
		} else {
			notification := model.Notification{
				Title:      keys["title"],
				Body:       keys["body"],
				SendTime:   time.Now(),
				StatusCode: status.StatusCode,
			}

			notificationID, err := app.NotificationRepository.Insert(&notification)

			if err != nil {
				responseStatus = ResponseStatusError
				data["err_code"] = 2
				data["err_text"] = "insert to db error"
				data["err_server"] = err
			} else {
				responseStatus = ResponseStatusOk
				data["notification_obj"] = notification
				data["notification_id"] = notificationID
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

func DeleteNotificationHistory(app *App, session *model.Session, route *Route, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := make(map[string]interface{})
	var responseStatus = ResponseStatusOk

	r.ParseMultipartForm(0)

	keys := make(map[string]string)
	keys["notification_id"] = r.FormValue("notification_id")

	emptyKeys := emptyValidator(keys)

	if len(emptyKeys) != 0 {
		responseStatus = ResponseStatusError
		data["err_code"] = 2
		data["items"] = emptyKeys
	} else {
		hexNotificationID, err := hex.DecodeString(keys["notification_id"])

		if err != nil || len(hexNotificationID) != 12 {
			responseStatus = ResponseStatusError
			data["err_code"] = 0
			data["err_text"] = "notification id format not valid"
		} else {
			notification, err := app.NotificationRepository.FindByID(bson.ObjectId(hexNotificationID))

			if err != nil {
				responseStatus = ResponseStatusError
				data["err_code"] = 0
				data["err_text"] = "notification not found"
			} else {
				err = app.NotificationRepository.Delete(notification)

				if err != nil {
					responseStatus = ResponseStatusError
					data["err_code"] = 0
					data["err_text"] = "cant delete notification!"
				} else {
					responseStatus = ResponseStatusOk
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
