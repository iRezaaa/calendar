package src

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"encoding/json"
	"gitlab.com/irezaa/calendar/src/model"
)

func GetAllIndicators(app *App, session *model.Session, route *Route, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := make(map[string]interface{})
	var responseStatus = ResponseStatusOk

	indicators, err := app.IndicatorRepository.FindAll()

	if err != nil {
		responseStatus = ResponseStatusError
		data["err_code"] = 0
	} else {
		responseStatus = ResponseStatusOk
		data["indicators"] = indicators
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

func AddIndicator(app *App, session *model.Session, route *Route, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := make(map[string]interface{})
	var responseStatus = ResponseStatusOk

	r.ParseMultipartForm(0)

	keys := make(map[string]string)
	keys["indicator_id"] = r.FormValue("indicator_id")
	keys["title"] = r.FormValue("title")
	keys["value"] = r.FormValue("value")
	keys["percent"] = r.FormValue("percent")

	emptyKeys := emptyValidator(keys)

	if len(emptyKeys) != 0 {
		responseStatus = ResponseStatusError
		data["err_code"] = 2
		data["items"] = emptyKeys
	} else {
		indicator := model.Indicator{
			ID:      keys["indicator_id"],
			Title:   keys["title"],
			Value:   keys["value"],
			Percent: keys["percent"],
		}

		changeInfo, err := app.IndicatorRepository.UpdateOrInsert(indicator)

		if err != nil {
			responseStatus = ResponseStatusError
			data["err_code"] = 0
		} else {
			responseStatus = ResponseStatusOk
			data["indicator_id"] = indicator.ID
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

func UpdateIndicator(app *App, session *model.Session, route *Route, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := make(map[string]interface{})
	var responseStatus = ResponseStatusOk

	r.ParseMultipartForm(0)

	keys := make(map[string]string)
	keys["indicator_id"] = r.FormValue("indicator_id")
	keys["title"] = r.FormValue("title")
	keys["value"] = r.FormValue("value")
	keys["percent"] = r.FormValue("percent")

	emptyKeys := emptyValidator(keys)

	if len(emptyKeys) != 0 {
		responseStatus = ResponseStatusError
		data["err_code"] = 2
		data["items"] = emptyKeys
	} else {
		indicator, err := app.IndicatorRepository.FindByID(keys["indicator_id"])

		if err != nil {
			responseStatus = ResponseStatusError
			data["err_code"] = 2
			data["err_text"] = "indicator not found!"
		} else {
			indicator.Title = keys["title"]
			indicator.Value = keys["value"]
			indicator.Percent = keys["percent"]
			_, err := app.IndicatorRepository.UpdateOrInsert(indicator)

			if err != nil {
				responseStatus = ResponseStatusError
				data["err_code"] = 0
				data["err_text"] = "error while update indicator"
			} else {
				responseStatus = ResponseStatusOk
				data["indicator_id"] = indicator.ID
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

func DeleteIndicator(app *App, session *model.Session, route *Route, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := make(map[string]interface{})
	var responseStatus = ResponseStatusOk

	r.ParseMultipartForm(0)

	keys := make(map[string]string)
	keys["indicator_id"] = r.FormValue("indicator_id")
	emptyKeys := emptyValidator(keys)

	if len(emptyKeys) != 0 {
		responseStatus = ResponseStatusError
		data["err_code"] = 2
		data["items"] = emptyKeys
	} else {
		indicator, err := app.IndicatorRepository.FindByID(keys["indicator_id"])

		if err != nil {
			responseStatus = ResponseStatusError
			data["err_code"] = 2
			data["err_text"] = "indicator not found!"
		} else {
			err := app.IndicatorRepository.Delete(indicator)

			if err != nil {
				responseStatus = ResponseStatusError
				data["err_code"] = 0
				data["err_text"] = "error while update indicator"
				data["err_server"] = err.Error()
			} else {
				responseStatus = ResponseStatusOk
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
