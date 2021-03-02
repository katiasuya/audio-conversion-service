package handler

import (
	"net/http"

	"github.com/katiasuya/audio-conversion-service/internal/web"
)

// ShowDoc shows service documentation.
func ShowDoc(w http.ResponseWriter, r *http.Request) {
	web.Respond(w, http.StatusOK, "Showing documentation")
}
