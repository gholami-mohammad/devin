package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"devin/routes"
)

func Init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	r := mux.NewRouter()

	u := r.PathPrefix("/api").Subrouter().StrictSlash(true)
	routes.LoadRoutes(u)

	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	srv := &http.Server{
		Handler:      r,
		Addr:         ":13000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
