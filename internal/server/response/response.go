package response

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/katiasuya/audio-conversion-service/internal/logging"
)

// Respond is a function to make http responses.
func Respond(w http.ResponseWriter, code int, payload interface{}) {
	logger := logging.Init().WithField("package", "response")

	w.Header().Set("Content-Type", "application/json")

	body, err := json.Marshal(payload)
	if err != nil {
		msg := fmt.Errorf("can't marshal the given payload: %w", err).Error()
		logger.Errorln(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	_, err = w.Write(body)
	if err != nil {
		msg := fmt.Errorf("can't write response: %w", err).Error()
		logger.Errorln(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(code)
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
