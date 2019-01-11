package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"fmt"
	"gitlab.com/irezaa/calendar/src/model"
	"time"
	"golang.org/x/crypto/bcrypt"
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
		Access:  AccessPublic,
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
		Access:  AccessPublic,
	},
	{
		Name:    "Get Day",
		Method:  MethodGet,
		Path:    "/day",
		Handler: GetDay,
		Access:  AccessPublic,
	},
	{
		Name:    "Grant User To Admin",
		Method:  MethodPost,
		Path:    "/admin/makeadmin",
		Access:  AccessAdmin,
		Handler: GrantUserToAdmin,
	},
	{
		Name:    "Get all banners",
		Method:  MethodGet,
		Path:    "/banner/all",
		Access:  AccessPublic,
		Handler: GetAllBanners,
	},
	{
		Name:    "Add banner",
		Method:  MethodPost,
		Path:    "/admin/banner/add",
		Access:  AccessAdmin,
		Handler: AddBanner,
	},
	{
		Name:    "Remove banner",
		Method:  MethodPost,
		Path:    "/admin/banner/delete",
		Access:  AccessAdmin,
		Handler: DeleteBanner,
	},
	{
		Name:    "Get Banner Picture",
		Method:  MethodGet,
		Path:    "/banner/picture",
		Access:  AccessPublic,
		Handler: GetBannerFile,
	},
	{
		Name:    "Get all indicators",
		Method:  MethodGet,
		Path:    "/indicator/all",
		Access:  AccessPublic,
		Handler: GetAllIndicators,
	},
	{
		Name:    "Add Indicator",
		Method:  MethodPost,
		Path:    "/admin/indicator/add",
		Access:  AccessAdmin,
		Handler: AddIndicator,
	},
	{
		Name:    "Update Indicator",
		Method:  MethodPost,
		Path:    "/admin/indicator/update",
		Access:  AccessAdmin,
		Handler: UpdateIndicator,
	},
	{
		Name:    "Delete Indicator",
		Method:  MethodPost,
		Path:    "/admin/indicator/delete",
		Access:  AccessAdmin,
		Handler: DeleteIndicator,
	},
	{
		Name:    "Get All News",
		Method:  MethodGet,
		Path:    "/news/all",
		Access:  AccessPublic,
		Handler: GetAllNews,
	},
	{
		Name:    "Add News",
		Method:  MethodPost,
		Path:    "/admin/news/add",
		Access:  AccessAdmin,
		Handler: AddNews,
	},
	{
		Name:    "Update News",
		Method:  MethodPost,
		Path:    "/admin/news/update",
		Access:  AccessAdmin,
		Handler: UpdateNews,
	},
	{
		Name:    "Delete News",
		Method:  MethodPost,
		Path:    "/admin/news/delete",
		Access:  AccessAdmin,
		Handler: DeleteNews,
	},
	{
		Name:    "Get Notifications History",
		Method:  MethodGet,
		Path:    "/admin/notification/all",
		Access:  AccessAdmin,
		Handler: GetNotificationHistory,
	},
	{
		Name:    "Push Notification to All Devices",
		Method:  MethodPost,
		Path:    "/admin/notification/push",
		Access:  AccessAdmin,
		Handler: PushNotification,
	},
	{
		Name:    "Delete Notification History",
		Method:  MethodPost,
		Path:    "/admin/notification/delete",
		Access:  AccessAdmin,
		Handler: DeleteNotificationHistory,
	},
}

func main() {

	app := new(App)
	app.Init(Config{
		DBConfig: DatabaseConfig{
			URL:      "mongodb:27017",
			DBName:   "calendar",
			UserName: "reza",
			Password: "!4P&4-/]G5YXdgX<",
			Timeout:  60 * time.Second,
		},
		FCMConfig: FCMConfig{
			ServerKey:    "AAAAPHnYoWI:APA91bGom3HGXqybz758HA-ewNbjXsHyCjbGTP_ZhC-JQ6pC2iynRunT7kcpswdrlC3EuwGgSUaiFvc-HsrVNQozp3GbHVMNTBnpzMNwZY6yuGkfpDUfh9c9UyBubXGbuoZUSgSW1hGcUHK0JfYRLuKyNF62xP9Zuw",
			DefaultTopic: "events",
		},
	})

	user, err := app.UserRepository.FindByID("admin")

	if err != nil {
		print("admin user is not exist!")

		print("Creating admin user...")

		password := []byte("tfbM2cjRpJx9wjnu")
		hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)

		if err != nil {
			panic(err)
		}

		user = &model.User{
			ID:    "admin",
			FName: "System Admin",
			LName: "No Last Name",
			Pass:  string(hashedPassword),
			Type:  model.UserTypeAdmin,
		}

		err = app.UserRepository.Insert(user)

		if err != nil {
			panic(err)
		}
	} else {
		print("admin is already exist!")
	}

	go startNotificationHandler(app, true)

	//err = startRouter(app, apiRoutes, 8080)
	//
	//if err != nil {
	//	panic(err)
	//}

	startCrawler(app)
}

func Index(app *App, session *model.Session, route *Route, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome! from route : "+route.Name)
}

func Hello(app *App, session *model.Session, route *Route, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name")+" from route : "+route.Name)
}
