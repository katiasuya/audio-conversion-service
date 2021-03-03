// Package web provides web response functions.
package web

import (
	"encoding/json"
	"log"
	"net/http"
)

// Respond is a function to make http responses.
func Respond(w http.ResponseWriter, code int, payload interface{}) {
	body, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "can't marshal the given payload", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(body)
}

// RespondErr is a function to make http error responses.
func RespondErr(w http.ResponseWriter, code int, err error) {
	type error struct {
		Code    int
		Message string
	}

	respErr := error{
		Code:    code,
		Message: err.Error(),
	}
	Respond(w, code, respErr)
}
