// Package handler provides http handlers for the service.
package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/katiasuya/audio-conversion-service/internal/web"
)

func validateName(name string) error {
	//	return errors.New("invalid name")
	return nil
}

func validateFormat(format string, targetFormat string) error {
	//return errors.New("invalid format")
	return nil
}

// Convert implements audio conversion.
func Convert(w http.ResponseWriter, r *http.Request) {
	reader, err := r.MultipartReader()
	if err != nil {
		web.RespondErr(w, http.StatusBadRequest, err)
		return
	}

	type conversionRequest struct {
		File         io.Reader `json:"file"`
		TargetFormat string    `json:"targetFormat"`
	}
	var cr conversionRequest

	var fileName string
	for part, err := reader.NextPart(); err != io.EOF; part, err = reader.NextPart() {
		if err != nil {
			web.RespondErr(w, http.StatusBadRequest, err)
			return
		}

		switch part.FormName() {
		case "File":
			cr.File = part
			fileName = part.FileName()
		case "TargetFormat":
			err := json.NewDecoder(part).Decode(&cr.TargetFormat)
			if err != nil {
				web.RespondErr(w, http.StatusBadRequest, err)
				return
			}
		default:
			web.RespondErr(w, http.StatusBadRequest, errors.New("Unasked data"))
			return
		}
	}

	if fileName == "" {
		web.RespondErr(w, http.StatusBadRequest, errors.New("file is missing"))
		return
	}

	if cr.TargetFormat == "" {
		web.RespondErr(w, http.StatusBadRequest, errors.New("format is missing"))
		return
	}

	fileParts := strings.Split(fileName, ".")
	name := fileParts[0]
	if err := validateName(name); err != nil {
		web.RespondErr(w, http.StatusBadRequest, err)
		return
	}

	format := fileParts[1]
	if err := validateFormat(format, cr.TargetFormat); err != nil {
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
