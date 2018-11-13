package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"encoding/json"
	"gitlab.com/irezaa/calendar/src/model"
	"strconv"
	"time"
	"encoding/hex"
	"gopkg.in/mgo.v2/bson"
)

func GetAllNews(app *App, session *model.Session, route *Route, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := make(map[string]interface{})
	var responseStatus = ResponseStatusOk

	news, err := app.NewsRepository.FindAll()

	if err != nil {
		responseStatus = ResponseStatusError
		data["err_code"] = 0
		data["err_server"] = err.Error()
		data["err_text"] = "error while fetch news from database!"
	} else {
		responseStatus = ResponseStatusOk
		data["news"] = news
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

func AddNews(app *App, session *model.Session, route *Route, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := make(map[string]interface{})
	var responseStatus = ResponseStatusOk

	r.Body = http.MaxBytesReader(w, r.Body, 10*MB)

	if err := r.ParseMultipartForm(0); err != nil {
		responseStatus = ResponseStatusError
		data["err_code"] = 0
		data["err_server"] = err.Error()
		data["err_text"] = "file size!"
	} else {
		keys := make(map[string]string)
		keys["title"] = r.FormValue("title")
		keys["content"] = r.FormValue("content")
		keys["contentType"] = r.FormValue("content_type")
		emptyKeys := emptyValidator(keys)

		filePath, receiveFileErr := FileUpload(r, "file", "news", []string{".png", ".jpeg", ".jpg"})

		if receiveFileErr != nil {
			responseStatus = ResponseStatusError
			data["err_code"] = 0
			data["err_text"] = "error on upload file"
			data["err_server"] = receiveFileErr.Error()
		} else {
			if len(emptyKeys) != 0 {
				responseStatus = ResponseStatusError
				data["err_code"] = 2
				data["items"] = emptyKeys
			} else {
				contentType, err := strconv.Atoi(keys["content_type"])

				if err != nil {
					responseStatus = ResponseStatusError
					data["err_code"] = 2
					data["err_text"] = "error while cast content type"
					data["err_server"] = err
				} else {
					news := model.News{
						Title:          keys["title"],
						Content:        keys["content"],
						ContentType:    model.NewsContentType(contentType),
						ImageURL:       filePath,
						CreateTime:     time.Now(),
						LastUpdateTime: time.Now(),
					}

					id, err := app.NewsRepository.Insert(&news)

					if err != nil {
						responseStatus = ResponseStatusError
						data["err_code"] = 0
						data["err_text"] = "error while add to database"
						data["err_server"] = err.Error()
					} else {
						responseStatus = ResponseStatusOk
						data["news_id"] = id
						data["news_obj"] = news
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

func UpdateNews(app *App, session *model.Session, route *Route, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := make(map[string]interface{})
	var responseStatus = ResponseStatusOk

	r.ParseMultipartForm(0)

	keys := make(map[string]string)
	keys["news_id"] = r.FormValue("news_id")
	keys["title"] = r.FormValue("title")
	keys["content"] = r.FormValue("content")
	keys["content_type"] = r.FormValue("content_type")

	requiredKeys := make(map[string]string)
	requiredKeys["news_id"] = keys["news_id"]
	emptyKeys := emptyValidator(requiredKeys)

	if len(emptyKeys) != 0 {
		responseStatus = ResponseStatusError
		data["err_code"] = 2
		data["items"] = emptyKeys
	} else {
		hexNewsID, err := hex.DecodeString(keys["news_id"])

		if err != nil {
			responseStatus = ResponseStatusError
			data["err_code"] = 0
			data["err_server"] = err
			data["err_text"] = "error while cast news id"
		} else {
			news, err := app.NewsRepository.FindByID(bson.ObjectId(hexNewsID))

			if err != nil {
				responseStatus = ResponseStatusError
				data["err_code"] = 2
				data["err_text"] = "news object not found!"
			} else {
				var updated = false

				if len(keys["title"]) > 0 && news.Title != keys["title"] {
					news.Title = keys["title"]
					updated = true
				}

				if len(keys["content"]) > 0 && news.Title != keys["content"] {
					news.Content = keys["content"]
					updated = true
				}

				contentType, err := strconv.Atoi(keys["content_type"])
				if err != nil {
					contentType = int(model.NewsContentTypeHTML)
				}

				if len(keys["content_type"]) > 0 && strconv.Itoa(contentType) != keys["content_type"] {
					news.ContentType = model.NewsContentType(contentType)
					updated = true
				}

				if !updated {
					responseStatus = ResponseStatusError
					data["err_code"] = 0
					data["err_text"] = "0 field changed!"
				} else {
					err := app.NewsRepository.Update(&news)

					if err != nil {
						responseStatus = ResponseStatusError
						data["err_code"] = 0
						data["err_text"] = "error while update news object"
					} else {
						responseStatus = ResponseStatusOk
						data["news_obj"] = news
						data["news_id"] = news.ID
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

func DeleteNews(app *App, session *model.Session, route *Route, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := make(map[string]interface{})
	var responseStatus = ResponseStatusOk

	r.ParseMultipartForm(0)

	keys := make(map[string]string)
	keys["news_id"] = r.FormValue("news_id")

	emptyKeys := emptyValidator(keys)

	if len(emptyKeys) != 0 {
		responseStatus = ResponseStatusError
		data["err_code"] = 2
		data["items"] = emptyKeys
	} else {
		hexNewsID, err := hex.DecodeString(keys["news_id"])

		if err != nil || len(hexNewsID) != 12 {
			responseStatus = ResponseStatusError
			data["err_code"] = 0
			data["err_text"] = "news id format not valid"
		} else {
			news, err := app.NewsRepository.FindByID(bson.ObjectId(hexNewsID))

			if err != nil {
				responseStatus = ResponseStatusError
				data["err_code"] = 0
				data["err_text"] = "news not found"
			} else {
				err = app.NewsRepository.Delete(news)

				if err != nil {
					responseStatus = ResponseStatusError
					data["err_code"] = 0
					data["err_text"] = "cant delete news!"
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
