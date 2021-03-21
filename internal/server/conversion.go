package server

import (
	"context"
	"net/http"
	"os/exec"

	"github.com/google/uuid"
	"golang.org/x/sync/semaphore"
)

const maxRequests = 10

var sem = semaphore.NewWeighted(maxRequests)

func (s *Server) convert(w http.ResponseWriter, fileID, filename, sourceFormat, targetFormat, requestID string) {
	if err := sem.Acquire(context.Background(), 1); err != nil {
		RespondErr(w, http.StatusInternalServerError, err)
		return
	}

	if err := s.repo.UpdateStatus(requestID); err != nil {
		RespondErr(w, http.StatusInternalServerError, err)
		return
	}

	targetFileID, err := uuid.NewRandom()
	if err != nil {
		RespondErr(w, http.StatusInternalServerError, err)
		return
	}
	cmd := exec.Command("ffmpeg", "-i", s.storage.Path+"/"+fileID+"."+sourceFormat,
		s.storage.Path+"/"+targetFileID.String()+"."+targetFormat)
	if err := cmd.Run(); err != nil {
		RespondErr(w, http.StatusInternalServerError, err)
		return
	}

	targetID, err := s.repo.InsertAudio(filename, targetFormat, targetFileID.String())
	if err != nil {
		RespondErr(w, http.StatusInternalServerError, err)
		return
	}
	if err := s.repo.UpdateRequest(requestID, targetID); err != nil {
		RespondErr(w, http.StatusInternalServerError, err)
		return
	}

	sem.Release(1)
}
