package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"gitlab.com/irezaa/calendar/src/model"
	"os"
	"io"
	"errors"
	"path/filepath"
)

type HttpMethod int
type RouteHandler func(*App, *model.Session, *Route, http.ResponseWriter, *http.Request, httprouter.Params)
type AccessType int
type ResponseStatus string

const (
	MB = 1 << 20
)

const (
	AccessPublic AccessType = iota + 1
	AccessUser
	AccessAdmin
)

const (
	MethodPost HttpMethod = iota + 1
	MethodGet
	MethodPut
	MethodDelete
)

const (
	ResponseStatusOk    ResponseStatus = "ok"
	ResponseStatusError ResponseStatus = "err"
)

type HttpJsonResponse struct {
	ResponseStatus ResponseStatus         `json:"status"`
	Data           map[string]interface{} `json:"data"`
}

type Route struct {
	Name    string
	Method  HttpMethod
	Path    string
	Handler RouteHandler
	Access  AccessType
}

func startRouter(app *App, routes []Route, port int) error {
	router := httprouter.New()

	for _, route := range routes {
		switch route.Method {
		case MethodPost:
			router.POST(route.Path, accessMiddleware(app, route, route.Handler))
			break
		case MethodGet:
			router.GET(route.Path, accessMiddleware(app, route, route.Handler))
			break
		case MethodPut:
			router.PUT(route.Path, accessMiddleware(app, route, route.Handler))
			break
		case MethodDelete:
			router.DELETE(route.Path, accessMiddleware(app, route, route.Handler))
			break
		}
	}

	router.ServeFiles("/uploads/*filepath", http.Dir("/uploads/"))

	return http.ListenAndServe(":"+strconv.Itoa(port), router)
}

func accessMiddleware(app *App, route Route, handler RouteHandler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		switch route.Access {
		case AccessPublic:
			token := r.Header.Get("Authorization")

			if &token != nil {
				session, err := app.SessionRepository.FindByAuthToken(token)

				if session != nil && err == nil && (session.User.Type == model.UserTypeNormal || session.User.Type == model.UserTypeAdmin) {
					handler(app, session, &route, w, r, ps)
				} else {
					handler(app, nil, &route, w, r, ps)
				}
			} else {
				handler(app, nil, &route, w, r, ps)
			}
			break
		case AccessUser:
			token := r.Header.Get("Authorization")

			if &token != nil {
				session, err := app.SessionRepository.FindByAuthToken(token)

				if session != nil && err == nil && (session.User.Type == model.UserTypeNormal || session.User.Type == model.UserTypeAdmin) {
					handler(app, session, &route, w, r, ps)
				} else {
					http.Error(w, "Unauthorized", 401)
				}
			} else {
				http.Error(w, "Unauthorized", 401)
			}

			break
		case AccessAdmin:
			token := r.Header.Get("Authorization")

			if &token != nil {
				session, err := app.SessionRepository.FindByAuthToken(token)

				if session != nil && err == nil && session.User.Type == model.UserTypeAdmin {
					handler(app, session, &route, w, r, ps)
				} else {
					http.Error(w, "Unauthorized", 401)
				}
			} else {
				http.Error(w, "Unauthorized", 401)
			}
			break
		}
	}
}

func FileUpload(r *http.Request, requestKey string, pathToSave string, allowedMimeTypes []string) (string, error) {

	requestFile, handler, err := r.FormFile(requestKey)
	defer requestFile.Close()

	if err != nil {
		return "", err
	}

	mimeType := filepath.Ext(handler.Filename)

	if allowedMimeTypes != nil && len(allowedMimeTypes) > 0 {
		isAccepted := false

		for _, element := range allowedMimeTypes {
			if mimeType == element {
				isAccepted = true
				break
			}
		}

		if !isAccepted {
			return "", errors.New("mime type not accepted")
		}
	}

	desFilePath := "/uploads/" + pathToSave + "/" + handler.Filename

	desFile, err := os.OpenFile(desFilePath, os.O_WRONLY|os.O_CREATE, 0666)
	defer desFile.Close()

	if err != nil {
		return "", err
	}

	io.Copy(desFile, requestFile)

	return desFilePath, err
}
