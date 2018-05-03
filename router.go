package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"gitlab.com/irezaa/calendar/model"
)

type HttpMethod int
type RouteHandler func(*App,*model.Session, *Route, http.ResponseWriter, *http.Request, httprouter.Params)
type AccessType int
type ResponseStatus string

const (
	AccessPublic AccessType = iota + 1
	AccessUser
	AccessAdmin
)

const (
	MethodPost   HttpMethod = iota + 1
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

func startRouter(app *App , routes []Route, port int) error {
	router := httprouter.New()

	for _, route := range routes {
		switch route.Method {
		case MethodPost:
			router.POST(route.Path, accessMiddleware(app,route, route.Handler))
			break
		case MethodGet:
			router.GET(route.Path, accessMiddleware(app,route, route.Handler))
			break
		case MethodPut:
			router.PUT(route.Path, accessMiddleware(app,route, route.Handler))
			break
		case MethodDelete:
			router.DELETE(route.Path, accessMiddleware(app,route, route.Handler))
			break
		}
	}

	return http.ListenAndServe(":"+strconv.Itoa(port), router)
}

func accessMiddleware(app *App ,route Route, handler RouteHandler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		switch route.Access {
		case AccessPublic:
			handler(app,nil, &route, w, r, ps)
			break
		case AccessUser:
			token := r.Header.Get("Authorization")

			if &token != nil {
				session, err := app.SessionRepository.FindByAuthToken(token)

				if session != nil && err == nil && (session.User.Type == model.UserTypeNormal || session.User.Type == model.UserTypeAdmin) {
					handler(app,session, &route, w, r, ps)
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
					handler(app,session, &route, w, r, ps)
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
