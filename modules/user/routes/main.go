package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

func LoadRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from user's module"))
	})
	return r
}
