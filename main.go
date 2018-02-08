package main

import (
	"log"
	"net/http"

	"gogit/models"
)

func main() {
	var a models.City
	log.Println(a)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hey"))
	})
	http.ListenAndServe(":80", nil)
}
