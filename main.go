package main

import (
	"gogit/models"
	"log"
	"net/http"
)

func main() {
	var a models.City
	log.Println(a)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hey"))
	})
	http.ListenAndServe(":80", nil)
}
