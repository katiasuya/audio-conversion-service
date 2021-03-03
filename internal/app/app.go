// Package app provides the function to start the application and http handlers.
package app

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/katiasuya/audio-conversion-service/internal/handler"
)

func initRoutes(r *mux.Router) {
	r.HandleFunc("/docs", handler.ShowDoc).Methods("GET")
	r.HandleFunc("/user/signup", handler.SignUp).Methods("POST")
	r.HandleFunc("/user/login", handler.LogIn).Methods("POST")
	r.HandleFunc("/conversion", handler.Convert).Methods("POST")
	r.HandleFunc("/request_history", handler.ShowHistory).Methods("GET")
	r.HandleFunc("/download_audio/{id}", handler.Download).Methods("GET")
}

// Run starts running the application service
func Run() error {
	r := mux.NewRouter()
	initRoutes(r)

	return http.ListenAndServe(":8000", r)
}
