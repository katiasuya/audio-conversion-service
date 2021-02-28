package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type audio struct {
	Name   string `json:"name"`
	Format string `json:"format"`
	File   string `json:"file"`
}

type convertResponse struct {
	ID string `json:"id"`
}

func (a *audio) validate() error {
	// return  errors.New("invalid name/format")
	return nil
}

func audioExists(id string) error {
	// return  errors.New("there is no song with the given id")
	return nil
}

func handlerConvert(w http.ResponseWriter, r *http.Request) {
	var a audio
	err := json.NewDecoder(r.Body).Decode(&a)
	defer r.Body.Close()
	if err != nil {
		respondErr(w, http.StatusBadRequest, err)
		w.Write([]byte("aaaa"))
		return
	}

	if err := a.validate(); err != nil {
		respondErr(w, http.StatusBadRequest, err)
		return
	}

	convertResp := convertResponse{
		ID: "1fa85f64-5717-4562-b3fc-2c963f66afa5",
	}

	respond(w, http.StatusCreated, convertResp)
}

func handlerDownload(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := audioExists(id); err != nil {
		respondErr(w, http.StatusNotFound, err)
		return
	}

	audio := audio{
		Name:   "Euphoria",
		Format: "WAV",
		File:   "euphoria.wav",
	}
	fmt.Fprint(w, fmt.Sprint(audio))
}
