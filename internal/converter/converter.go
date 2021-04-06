// Package converter implements a function for audio converting.
package converter

import (
	"context"
	"os/exec"

	"github.com/google/uuid"
	"github.com/katiasuya/audio-conversion-service/internal/logger"
	"github.com/katiasuya/audio-conversion-service/internal/repository"
	"github.com/katiasuya/audio-conversion-service/internal/storage"
	"golang.org/x/sync/semaphore"
)

var status = []string{"processing", "done", "failed"}

// Converter converts audio files to other formats.
type Converter struct {
	sem  *semaphore.Weighted
	repo *repository.Repository
	st   *storage.Storage
}

// New creates a new Converter with given fields.
func New(sem *semaphore.Weighted, repo *repository.Repository, st *storage.Storage) *Converter {
	return &Converter{sem: sem, repo: repo, st: st}
}

// Convert implements audio conversion.
func (c *Converter) Convert(fileID, filename, sourceFormat, targetFormat, requestID string) {
	logger := logger.Init().WithField("package", "converter")

	if err := c.sem.Acquire(context.Background(), 1); err != nil {
		if err1 := c.repo.UpdateRequest(requestID, status[2], ""); err1 != nil {
			logger.Errorln(err1)
		}
		logger.Errorln("can't acquire the semaphore: ", err)
		return
	}

	if err := c.repo.UpdateRequest(requestID, status[0], ""); err != nil {
		if err1 := c.repo.UpdateRequest(requestID, status[2], ""); err1 != nil {
			logger.Errorln(err1)
		}
		logger.Errorln("can't update request: ", err)
		return
	}
	logger.Debugln("status changed to processing")

	targetFileID, err := uuid.NewRandom()
	if err != nil {
		if err1 := c.repo.UpdateRequest(requestID, status[2], ""); err1 != nil {
			logger.Errorln(err1)
		}
		logger.Errorln("can't generate targetFileID: ", err)
		return
	}
	logger.Debugln("targetFileID generated successfully")

	cmd := exec.Command("ffmpeg", "-i", c.st.Path+"/"+fileID+"."+sourceFormat,
		c.st.Path+"/"+targetFileID.String()+"."+targetFormat)
	if err := cmd.Run(); err != nil {
		if err1 := c.repo.UpdateRequest(requestID, status[2], ""); err1 != nil {
			logger.Errorln(err1)
		}
		logger.Errorln("can't convert the file: ", err)
		return
	}
	logger.Debugln("convertion performed successfully")

	targetID, err := c.repo.InsertAudio(filename, targetFormat, targetFileID.String())
	if err != nil {
		if err1 := c.repo.UpdateRequest(requestID, status[2], ""); err1 != nil {
			logger.Errorln(err1)
		}
		logger.Errorln("can't insert audio: ", err)
		return
	}
	logger.Debugln("converted audio inserted successfully")

	if err := c.repo.UpdateRequest(requestID, status[1], targetID); err != nil {
		if err1 := c.repo.UpdateRequest(requestID, status[2], ""); err1 != nil {
			logger.Errorln(err1)
		}
		logger.Errorln("can't update request: ", err)
		return
	}
	logger.Debugln("status changed to done")

	logger.WithField("fileID", fileID).Infoln("the file was converted successfully")
	c.sem.Release(1)
}
