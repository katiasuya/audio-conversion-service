// Package app provides the function to start the application and http handlers.
package app

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func initRoutes(r *mux.Router) {
	r.HandleFunc("/docs", showDoc).Methods("GET")
	r.HandleFunc("/user/signup", signUp).Methods("POST")
	r.HandleFunc("/user/login", logIn).Methods("POST")
	r.HandleFunc("/conversion", convert).Methods("POST")
	r.HandleFunc("/request_history", showHistory).Methods("GET")
	r.HandleFunc("/download_audio/{id}", download).Methods("GET")
}

// Start starts the application service
func Start() {
	r := mux.NewRouter()
	initRoutes(r)
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8000", nil))
}
