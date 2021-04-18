package response

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/katiasuya/audio-conversion-service/internal/logging"
)

// Respond is a function to make http responses.
func Respond(w http.ResponseWriter, code int, payload interface{}) {
	body, err := json.Marshal(payload)
	if err != nil {
		msg := fmt.Errorf("can't marshal the given payload: %w", err).Error()
		logging.Init().WithField("package", "response").Errorln(msg)
		http.Error(w, msg, http.StatusInternalServerError)
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
