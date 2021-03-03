// Package handler provides http handlers for the service.
package handler

import (
	"crypto/rand"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/katiasuya/audio-conversion-service/internal/web"
)

func audioExists(id string) error {
	//return errors.New("there is no song with the given id")
	return nil
}

// Download implements audio downloading.
func Download(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := audioExists(id); err != nil {
		web.RespondErr(w, http.StatusNotFound, err)
		return
	}

	type audio struct {
		Name   string `json:"name"`
		Format string `json:"format"`
	}
	file := audio{
		Name:   "Euphoria.mp3",
		Format: "mp3",
	}

	mwriter := multipart.NewWriter(w)
	w.Header().Set("Content-Type", mwriter.FormDataContentType())

	fileParts := []string{"nameFormat", "file"}
	for i := 0; i < len(fileParts); i++ {
		fw, err := mwriter.CreateFormField(fileParts[i])
		if err != nil {
			web.RespondErr(w, http.StatusInternalServerError, err)
			return
		}

		switch fileParts[i] {
		case "nameFormat":
			if err := json.NewEncoder(fw).Encode(&file); err != nil {
				web.RespondErr(w, http.StatusInternalServerError, err)
				return
			}
		case "file":
			if _, err := io.CopyN(fw, rand.Reader, 32); err != nil {
				web.RespondErr(w, http.StatusInternalServerError, err)
				return
			}
		}
	}

	if err := mwriter.Close(); err != nil {
		web.RespondErr(w, http.StatusInternalServerError, err)
		return
	}
}
