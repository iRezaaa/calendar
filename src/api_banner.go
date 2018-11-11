package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"encoding/json"
	"gitlab.com/irezaa/calendar/src/model"
	"encoding/hex"
	"gopkg.in/mgo.v2/bson"
	"os"
	"io"
	"path/filepath"
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
		data["err_server"] = err.Error()
		data["err_text"] = "file size!"
	} else {
		keys := make(map[string]string)
		keys["title"] = r.FormValue("title")
		keys["link"] = r.FormValue("link")
		emptyKeys := emptyValidator(keys)

		filePath, receiveFileErr := FileUpload(r, "file", "banners", []string{".png", ".jpeg",".jpg"})

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

func GetBannerFile(app *App, session *model.Session, route *Route, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := make(map[string]interface{})
	var responseStatus = ResponseStatusOk
	var responseImage *os.File
	var contentType = ""

	q := r.URL.Query()
	keys := make(map[string]string)
	keys["banner_id"] = q.Get("banner_id")
	emptyKeys := emptyValidator(keys)

	if len(emptyKeys) != 0 {
		responseStatus = ResponseStatusError
		data["err_code"] = 0
		data["items"] = emptyKeys
		data["err_text"] = "missing required items!"
	}else{
		hexBannerID, err := hex.DecodeString(keys["banner_id"])

		if err != nil || len(hexBannerID) != 12 {
			responseStatus = ResponseStatusError
			data["err_code"] = 0
			data["err_text"] = "banner id format not valid"
		} else {
			banner , err := app.BannerRepository.FindByID(bson.ObjectId(hexBannerID))

			if err != nil {
				responseStatus = ResponseStatusError
				data["err_code"] = 0
				data["err_text"] = "banner with given id not found!"
			}else{
				databaseImageAddress := banner.ImageURL

				if databaseImageAddress == "" {
					responseStatus = ResponseStatusError
					data["err_code"] = 0
					data["err_text"] = "banner address in database is not valid!"
				}else{
					responseImage, err = os.Open("/uploads/" + databaseImageAddress)

					if err != nil {
						responseStatus = ResponseStatusError
						data["err_code"] = 0
						data["err_text"] = "error while open file from db!"
					}else {
						responseStatus = ResponseStatusOk
						print("/uploads/" + databaseImageAddress)
					}
				}
			}
		}

	}

	if responseImage != nil {
		ext := filepath.Ext(responseImage.Name())

		if ext == "" {
			responseStatus = ResponseStatusError
			data["err_code"] = 0
			data["err_text"] = "cannot get file type!!"
		}else {
			if ext == ".png" {
				contentType = "image/png"
			}else if ext == ".jpg" || ext == ".jpeg" {
				contentType = "image/jpeg"
			}else {
				responseStatus = ResponseStatusError
				data["err_code"] = 0
				data["err_text"] = "cannot get file type!!"
			}
		}

		print("ext := " + ext)
	}


	if responseStatus == ResponseStatusError {
		response := HttpJsonResponse{
			ResponseStatus: responseStatus,
			Data:           data,
		}

		jsonResponse, err := json.Marshal(response)

		if err == nil {
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonResponse)
		}
	}else if responseImage != nil && contentType != "" {
		defer responseImage.Close()
		w.Header().Set("Content-Type", contentType) // <-- set the content-type header
		print("\n banner image writing to response!")
		io.Copy(w, responseImage)
	}

	print("\n exit!!!!!!!!!!!!!!kasjdkajdkasjd")

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
