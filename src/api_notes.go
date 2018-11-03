package src

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"encoding/json"
	"gitlab.com/irezaa/calendar/src/model"
	"time"
	"strconv"
	"encoding/hex"
	"gopkg.in/mgo.v2/bson"
	"github.com/yaa110/go-persian-calendar/ptime"
)

func GetNotes(app *App, session *model.Session, route *Route, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := make(map[string]interface{})
	var responseStatus = ResponseStatusOk

	notes, err := app.NoteRepository.FindByUserID(session.User.ID)

	if err != nil {
		responseStatus = ResponseStatusError
		data["err_code"] = 0
	} else {
		responseStatus = ResponseStatusOk
		data["notes"] = notes
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

func AddNote(app *App, session *model.Session, route *Route, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := make(map[string]interface{})
	var responseStatus = ResponseStatusOk

	r.ParseMultipartForm(0)

	keys := make(map[string]string)
	keys["title"] = r.FormValue("title")
	keys["body"] = r.FormValue("body")
	keys["show_timeStamp"] = r.FormValue("show_timeStamp")

	emptyKeys := emptyValidator(keys)

	if len(emptyKeys) != 0 {
		responseStatus = ResponseStatusError
		data["err_code"] = 2
		data["items"] = emptyKeys
	} else {
		responseStatus = ResponseStatusOk

		var showTimestampInMills int64

		showTimestampInMills = 0

		showTimestampInMills, _ = strconv.ParseInt(keys["show_timeStamp"], 10, 64)

		gregorianTime := time.Unix(showTimestampInMills, 0).In(ptime.Iran())

		note := model.Note{
			Title:  keys["title"],
			Body:   keys["body"],
			UserID: session.User.ID,
			Day:    gregorianTime.Day(),
			Month:  int(gregorianTime.Month()),
			Year:   int(gregorianTime.Year()),
		}

		changeInfo, err := app.NoteRepository.UpdateOrInsert(note)

		if err != nil {
			responseStatus = ResponseStatusError
			data["err_code"] = 0
		} else {
			responseStatus = ResponseStatusOk
			data["note_id"] = changeInfo.UpsertedId
			data["updated"] = changeInfo.Updated
			data["matched"] = changeInfo.Matched
			data["removed"] = changeInfo.Removed
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

func UpdateNote(app *App, session *model.Session, route *Route, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := make(map[string]interface{})
	var responseStatus = ResponseStatusOk

	r.ParseMultipartForm(0)

	keys := make(map[string]string)
	keys["note_id"] = r.FormValue("note_id")

	emptyKeys := emptyValidator(keys)

	keys["body"] = r.FormValue("body")
	keys["title"] = r.FormValue("title")

	if len(emptyKeys) != 0 {
		responseStatus = ResponseStatusError
		data["err_code"] = 2
		data["items"] = emptyKeys
	} else {
		hexNoteID, err := hex.DecodeString(keys["note_id"])

		if err != nil || len(hexNoteID) != 12 {
			responseStatus = ResponseStatusError
			data["err_code"] = 0
			data["err_text"] = "note id format not valid"
		} else {
			note, err := app.NoteRepository.FindByID(bson.ObjectId(hexNoteID))

			if err != nil {
				responseStatus = ResponseStatusError
				data["err_code"] = 0
				data["err_text"] = "note not found"
			} else {
				if note.UserID == session.User.ID {
					if keys["title"] != "" {
						note.Title = keys["title"]
					}

					if keys["body"] != "" {
						note.Body = keys["body"]
					}

					changeInfo, err := app.NoteRepository.UpdateOrInsert(note)

					if err != nil {
						responseStatus = ResponseStatusError
						data["err_code"] = 0
						data["err_text"] = "cant update note!"
					} else {
						responseStatus = ResponseStatusOk
						data["note_id"] = note.ID
						data["updated"] = changeInfo.Updated
						data["matched"] = changeInfo.Matched
						data["removed"] = changeInfo.Removed
					}
				} else {
					responseStatus = ResponseStatusError
					data["err_code"] = 1
					data["err_text"] = "you are not note owner :D"
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

func DeleteNote(app *App, session *model.Session, route *Route, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := make(map[string]interface{})
	var responseStatus = ResponseStatusOk

	r.ParseMultipartForm(0)

	keys := make(map[string]string)
	keys["note_id"] = r.FormValue("note_id")

	emptyKeys := emptyValidator(keys)

	if len(emptyKeys) != 0 {
		responseStatus = ResponseStatusError
		data["err_code"] = 2
		data["items"] = emptyKeys
	} else {
		hexNoteID, err := hex.DecodeString(keys["note_id"])

		if err != nil || len(hexNoteID) != 12 {
			responseStatus = ResponseStatusError
			data["err_code"] = 0
			data["err_text"] = "note id format not valid"
		} else {
			note, err := app.NoteRepository.FindByID(bson.ObjectId(hexNoteID))

			if err != nil {
				responseStatus = ResponseStatusError
				data["err_code"] = 0
				data["err_text"] = "note not found"
			} else {
				if note.UserID == session.User.ID {
					err = app.NoteRepository.Delete(note)

					if err != nil {
						responseStatus = ResponseStatusError
						data["err_code"] = 0
						data["err_text"] = "cant delete note!"
					} else {
						responseStatus = ResponseStatusOk
					}
				} else {
					responseStatus = ResponseStatusError
					data["err_code"] = 1
					data["err_text"] = "you are not note owner :D"
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
