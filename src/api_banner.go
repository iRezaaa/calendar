package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"encoding/json"
	"gitlab.com/irezaa/calendar/src/model"
	"encoding/hex"
	"gopkg.in/mgo.v2/bson"
)

func GetAllBanners(app *App, session *model.Session, route *Route, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := make(map[string]interface{})
	var responseStatus = ResponseStatusOk

	banners, err := app.BannerRepository.FindAll()

	if err != nil {
		responseStatus = ResponseStatusError
		data["err_code"] = 0
	} else {
		responseStatus = ResponseStatusOk
		data["banners"] = banners
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

func AddBanner(app *App, session *model.Session, route *Route, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := make(map[string]interface{})
	var responseStatus = ResponseStatusOk

	r.Body = http.MaxBytesReader(w, r.Body, 10*MB)
	if err := r.ParseMultipartForm(10 * MB); err != nil {
		responseStatus = ResponseStatusError
		data["err_code"] = 0
		data["err_text"] = "file size!"
	} else {
		keys := make(map[string]string)
		keys["title"] = r.FormValue("title")
		keys["link"] = r.FormValue("link")
		emptyKeys := emptyValidator(keys)

		filePath, receiveFileErr := FileUpload(r, "file", "banners", []string{"image/png", "image/jpeg"})

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
				banner := model.Banner{
					Title:    keys["title"],
					Link:     keys["link"],
					ImageURL: filePath,
				}

				_, err := app.BannerRepository.UpdateOrInsert(banner)

				if err != nil {
					responseStatus = ResponseStatusError
					data["err_code"] = 0
					data["err_text"] = "error while add to database"
					data["err_server"] = err.Error()
				} else {
					responseStatus = ResponseStatusOk
					data["banner"] = banner
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

func DeleteBanner(app *App, session *model.Session, route *Route, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := make(map[string]interface{})
	var responseStatus = ResponseStatusOk

	r.ParseMultipartForm(0)

	keys := make(map[string]string)
	keys["banner_id"] = r.FormValue("banner_id")

	emptyKeys := emptyValidator(keys)

	if len(emptyKeys) != 0 {
		responseStatus = ResponseStatusError
		data["err_code"] = 2
		data["items"] = emptyKeys
	} else {
		hexBannerID, err := hex.DecodeString(keys["banner_id"])

		if err != nil || len(hexBannerID) != 12 {
			responseStatus = ResponseStatusError
			data["err_code"] = 0
			data["err_text"] = "banner id format not valid"
		} else {
			banner, err := app.BannerRepository.FindByID(bson.ObjectId(hexBannerID))

			if err != nil {
				responseStatus = ResponseStatusError
				data["err_code"] = 0
				data["err_text"] = "banner not found"
			} else {
				err = app.BannerRepository.Delete(banner)

				if err != nil {
					responseStatus = ResponseStatusError
					data["err_code"] = 0
					data["err_text"] = "cant delete banner!"
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
