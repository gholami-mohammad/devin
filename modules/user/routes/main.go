package routes

import (
	"github.com/gorilla/mux"

	user_ctrl "devin/modules/user/controllers"
)

func LoadRoutes(r *mux.Router) *mux.Router {
	r.HandleFunc("/signup", user_ctrl.Signup)
	r.HandleFunc("/signin", user_ctrl.Signin)

	return r
}
