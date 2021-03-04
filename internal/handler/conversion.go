// Package handler provides http handlers for the service.
package handler

import (
	"errors"
	"net/http"

	"github.com/katiasuya/audio-conversion-service/internal/web"
)

func validateName(name string) error {
	//	return errors.New("invalid name")
	return nil
}

func validateFormats(sourceFormat string, targetFormat string) error {
	//return errors.New("invalid source/target format")
	return nil
}

// Convert implements audio conversion.
func Convert(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		web.RespondErr(w, http.StatusBadRequest, err)
		return
	}
	defer file.Close()

	sourceFormat := r.FormValue("sourceFormat")
	targetFormat := r.FormValue("targetFormat")
	if sourceFormat == "" || targetFormat == "" {
		web.RespondErr(w, http.StatusBadRequest, errors.New("source/target format is missing"))
		return
	}

	name := header.Filename
	if err := validateName(name); err != nil {
		web.RespondErr(w, http.StatusBadRequest, err)
		return
	}

	if err := validateFormats(sourceFormat, targetFormat); err != nil {
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
