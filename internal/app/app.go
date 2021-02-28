// Package app provides the function to start the application and http handlers.
package app

import (
	"net/http"

	"github.com/gorilla/mux"
)

func initRoutes(r *mux.Router) {
	r.HandleFunc("/docs", handlerShowDoc).Methods("GET")
	r.HandleFunc("/user/signup", handlerSignUp).Methods("POST")
	r.HandleFunc("/user/login", handlerLogIn).Methods("POST")
	r.HandleFunc("/conversion", handlerConvert).Methods("POST")
	r.HandleFunc("/request_history", handlerShowHistory).Methods("GET")
	r.HandleFunc("/download_audio/{id}", handlerDownload).Methods("GET")
}

// Run starts running the application service
func Run() error {
	r := mux.NewRouter()
	initRoutes(r)

	return http.ListenAndServe(":8000", r)
}
