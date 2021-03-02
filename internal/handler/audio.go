// Package handler provides http handlers for the service.
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/katiasuya/audio-conversion-service/internal/web"
)

type audio struct {
	Name   string `json:"name"`
	Format string `json:"format"`
	File   string `json:"file"`
}

func (a *audio) validate() error {
	// return  errors.New("invalid name/format")
	return nil
}

func audioExists(id string) error {
	// return  errors.New("there is no song with the given id")
	return nil
}

// Convert implements audio conversion.
func Convert(w http.ResponseWriter, r *http.Request) {
	var a audio
	err := json.NewDecoder(r.Body).Decode(&a)
	defer r.Body.Close()
	if err != nil {
		web.RespondErr(w, http.StatusBadRequest, err)
		return
	}

	if err := a.validate(); err != nil {
		web.RespondErr(w, http.StatusBadRequest, err)
		return
	}

	type response struct {
		ID string `json:"id"`
	}
	convertResp := response{
		ID: "1fa85f64-5717-4562-b3fc-2c963f66afa5",
	}
	web.Respond(w, http.StatusCreated, convertResp)
}

// Download implements audio downloading.
func Download(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := audioExists(id); err != nil {
		web.RespondErr(w, http.StatusNotFound, err)
		return
	}

	audioResp := audio{
		Name:   "Euphoria",
		Format: "WAV",
		File:   "euphoria.wav",
	}
	web.Respond(w, http.StatusOK, audioResp)
}
