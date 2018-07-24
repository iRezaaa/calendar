package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"fmt"
	"gitlab.com/irezaa/calendar/model"
)

var apiRoutes = []Route{
	{
		Name:    "Home Page",
		Method:  MethodGet,
		Path:    "/",
		Handler: Index,
		Access:  AccessPublic,
	},
	{
		Name:    "Login",
		Method:  MethodPost,
		Path:    "/login",
		Handler: Login,
		Access:  AccessPublic,
	},
	{
		Name:    "Register",
		Method:  MethodPost,
		Path:    "/register",
		Handler: Register,
		Access:  AccessPublic,
	},
	{
		Name:    "Get User Info",
		Method:  MethodGet,
		Path:    "/user_info",
		Handler: GetUserInfo,
		Access:  AccessUser,
	},
	{
		Name:    "Get User Notes",
		Method:  MethodGet,
		Path:    "/note/all",
		Handler: GetNotes,
		Access:  AccessUser,
	},
	{
		Name:    "Add Note",
		Method:  MethodPost,
		Path:    "/note/add",
		Handler: AddNote,
		Access:  AccessUser,
	}, {
		Name:    "Update Note",
		Method:  MethodPost,
		Path:    "/note/update",
		Handler: UpdateNote,
		Access:  AccessUser,
	}, {
		Name:    "Delete Note",
		Method:  MethodPost,
		Path:    "/note/delete",
		Handler: DeleteNote,
		Access:  AccessUser,
	}, {
		Name:    "Get Month",
		Method:  MethodGet,
		Path:    "/events/all",
		Handler: GetEvents,
		Access:  AccessUser,
	}, {
		Name:    "Add Personal Event",
		Method:  MethodPost,
		Path:    "/events/add",
		Handler: AddPersonalEvent,
		Access:  AccessUser,
	}, {
		Name:    "Delete Personal Event",
		Method:  MethodPost,
		Path:    "/events/delete",
		Handler: DeletePersonalEvent,
		Access:  AccessUser,
	}, {
		Name:    "Add Panel Event",
		Method:  MethodPost,
		Path:    "/admin/events/add",
		Handler: AddPanelEvent,
		Access:  AccessAdmin,
	}, {
		Name:    "Delete Panel Event",
		Method:  MethodPost,
		Path:    "/admin/events/delete",
		Handler: DeletePanelEvent,
		Access:  AccessAdmin,
	}, {
		Name:    "Get today",
		Method:  MethodGet,
		Path:    "/today",
		Handler: Today,
		Access:  AccessUser,
	},
}

func main() {

	app := new(App)
	app.Init(Config{
		MongoURL: "127.0.0.1:27017",
	})

	go startNotificationHandler(app,true)

	err := startRouter(app, apiRoutes, 8080)

	if err != nil {
		panic(err)
	}

	//startCrawler(app)
}

func Index(app *App, session *model.Session, route *Route, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome! from route : "+route.Name)
}

func Hello(app *App, session *model.Session, route *Route, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name")+" from route : "+route.Name)
}
