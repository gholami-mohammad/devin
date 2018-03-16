package main

import (
	"log"
	"net/http"
	"time"

	user_router "devin/modules/user/routes"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	u := r.PathPrefix("/api").Subrouter().StrictSlash(true)
	user_router.LoadRoutes(u)

	srv := &http.Server{
		Handler:      r,
		Addr:         ":13000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
