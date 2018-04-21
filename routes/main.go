package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"devin/middlewares"
	org_ctrl "devin/modules/organization/controllers"
	user_ctrl "devin/modules/user/controllers"
)

func LoadRoutes(r *mux.Router) *mux.Router {
	r.HandleFunc("/signup", user_ctrl.Signup).Methods(http.MethodPost)
	r.HandleFunc("/signin", user_ctrl.Signin).Methods(http.MethodPost)

	secureArea := r.PathPrefix("/").Subrouter().StrictSlash(true)
	secureArea.Use(middlewares.Authenticate)
	secureArea.HandleFunc("/user/{id:[0-9]+}/update", user_ctrl.UpdateProfile).Methods(http.MethodPost)
	secureArea.HandleFunc("/user/{id:[0-9]+}/update_username", user_ctrl.UpdateUsername).Methods(http.MethodPost)
	secureArea.HandleFunc("/user/{id:[0-9]+}/update_avatar", user_ctrl.UpdateAvatar).Methods(http.MethodPost)
	secureArea.HandleFunc("/user/{id:[0-9]+}/update_password", user_ctrl.UpdatePassword).Methods(http.MethodPost)
	secureArea.HandleFunc("/user/{id:[0-9]+}/organization/save", org_ctrl.Save).Methods(http.MethodPost)
	secureArea.HandleFunc("/whoami", user_ctrl.Whoami).Methods(http.MethodGet)
	secureArea.HandleFunc("/whois/{id:[0-9]+}", user_ctrl.Whois).Methods(http.MethodGet)
	secureArea.HandleFunc("/profile_basic_info", user_ctrl.ProfileBasicInfo).Methods(http.MethodGet)

	return r
}
