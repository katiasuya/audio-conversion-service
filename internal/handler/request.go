package handler

import (
	"net/http"

	"github.com/katiasuya/audio-conversion-service/internal/web"
)

// ShowHistory shows request history.
func ShowHistory(w http.ResponseWriter, r *http.Request) {
	type historyResponse struct {
		ID             string `json:"ID"`
		OriginalID     string `json:"originalID"`
		OriginalFormat string `json:"originalFormat"`
		TargetID       string `json:"targetID"`
		TargetFormat   string `json:"targetFormat"`
		Created        string `json:"created"`
		Updated        string `json:"updated"`
		Status         string `json:"status"`
	}
	resp := []historyResponse{
		{
			ID:             "1fa85f64-5717-4562-b3fc-2c963f66afa5",
			OriginalID:     "2fa85f64-5717-4562-b3fc-2c963f66afa5",
			OriginalFormat: "MP3",
			TargetID:       "3fa85f64-5717-4562-b3fc-2c963f66afa5",
			TargetFormat:   "WAV",
			Created:        "2020-02-20 T 11:32:28 Z",
			Updated:        "2020-02-20 T 12:32:28 Z",
			Status:         "done",
		},
	}
	web.Respond(w, http.StatusOK, resp)
}
