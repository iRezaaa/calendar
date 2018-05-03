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
		Access:  AccessAdmin,
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
}

func main() {

	app := new(App)
	app.Init(Config{
		MongoURL: "127.0.0.1:27017",
	})

	err := startRouter(app, apiRoutes, 8080)

	if err != nil {
		panic(err)
	}
}

func Index(app *App, session *model.Session, route *Route, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome! from route : "+route.Name)
}

func Hello(app *App, session *model.Session, route *Route, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name")+" from route : "+route.Name)
}
