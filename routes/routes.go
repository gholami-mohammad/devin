package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"devin/middlewares"
	org_ctrl "devin/modules/organization/controllers"
	project_ctrl "devin/modules/project/controllers"
	user_ctrl "devin/modules/user/controllers"
)

func LoadRoutes(r *mux.Router) *mux.Router {
	r.HandleFunc("/signup", user_ctrl.Signup).Methods(http.MethodPost)
	r.HandleFunc("/signup/verify", user_ctrl.VerifySignup).Methods(http.MethodGet)
	r.HandleFunc("/signin", user_ctrl.Signin).Methods(http.MethodPost)
	r.HandleFunc("/password_reset/request", user_ctrl.RequestPasswordReset).Methods(http.MethodPost)
	r.HandleFunc("/password_reset/validate", user_ctrl.ValidatePasswordResetLink).Methods(http.MethodGet)
	r.HandleFunc("/password_reset/do", user_ctrl.ResetPassword).Methods(http.MethodPost)

	secureArea := r.PathPrefix("/").Subrouter().StrictSlash(true)
	secureArea.Use(middlewares.Authenticate)
	secureArea.HandleFunc("/user/{id:[0-9]+}/update", user_ctrl.UpdateProfile).Methods(http.MethodPost)
	secureArea.HandleFunc("/user/{id:[0-9]+}/update_username", user_ctrl.UpdateUsername).Methods(http.MethodPost)
	secureArea.HandleFunc("/user/{id:[0-9]+}/update_email", user_ctrl.UpdateEmail).Methods(http.MethodPost)
	secureArea.HandleFunc("/user/{id:[0-9]+}/update_avatar", user_ctrl.UpdateAvatar).Methods(http.MethodPost)
	secureArea.HandleFunc("/user/{id:[0-9]+}/update_password", user_ctrl.UpdatePassword).Methods(http.MethodPost)
	secureArea.HandleFunc("/user/{id:[0-9]+}/organization/save", org_ctrl.Save).Methods(http.MethodPost)
	secureArea.HandleFunc("/user/{id:[0-9]+}/pending_invitations", org_ctrl.PendingInvitationRequests).Methods(http.MethodGet)

	secureArea.HandleFunc("/organization/list", org_ctrl.UserOrganizationsIndex).Methods(http.MethodGet)
	secureArea.HandleFunc("/organization/{id:[0-9]+}/invite_user", org_ctrl.InviteUser).Methods(http.MethodPost)
	secureArea.HandleFunc("/organization/{organization_id:[0-9]+}/user/{user_id:[0-9]+}/update_permissions", org_ctrl.UpdateUserPermissionsOnOrganization).Methods(http.MethodPost)

	secureArea.HandleFunc("/invitation/{id:[0-9]+}/set_acceptance/{acceptance_status:(?:accept|reject)}", org_ctrl.AcceptOrRejectInvitation)

	secureArea.HandleFunc("/projects", project_ctrl.ProjectController{}.ProjectsIndex).Methods(http.MethodGet)
	secureArea.HandleFunc("/projects/basic_info", project_ctrl.ProjectController{}.BasicInfo)

	secureArea.HandleFunc("/whoami", user_ctrl.Whoami).Methods(http.MethodGet)
	secureArea.HandleFunc("/whois/{id:[0-9]+}", user_ctrl.Whois).Methods(http.MethodGet)
	secureArea.HandleFunc("/profile_basic_info", user_ctrl.ProfileBasicInfo).Methods(http.MethodGet)

	return r
}
