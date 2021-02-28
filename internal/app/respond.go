package app

import (
	"encoding/json"
	"log"
	"net/http"
)

// Error represents any error appeared during the running service
type Error struct {
	Code    int
	Message string
}

func respond(w http.ResponseWriter, code int, payload interface{}) {
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

func respondErr(w http.ResponseWriter, code int, err error) {
	respErr := Error{
		Code:    code,
		Message: err.Error(),
	}
	respond(w, code, respErr)
}
