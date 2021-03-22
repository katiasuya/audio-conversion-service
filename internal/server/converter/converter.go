// Package converter implements a function for audio converting.
package converter

import (
	"context"
	"log"
	"os/exec"

	"github.com/google/uuid"
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
	if err := c.sem.Acquire(context.Background(), 1); err != nil {
		if err1 := c.repo.UpdateRequest(requestID, status[2], ""); err1 != nil {
			log.Println(err1)
		}
		log.Println(err)
		return
	}

	if err := c.repo.UpdateRequest(requestID, status[0], ""); err != nil {
		if err1 := c.repo.UpdateRequest(requestID, status[2], ""); err1 != nil {
			log.Println(err1)
		}
		log.Println(err)
		return
	}

	targetFileID, err := uuid.NewRandom()
	if err != nil {
		if err1 := c.repo.UpdateRequest(requestID, status[2], ""); err1 != nil {
			log.Println(err1)
		}
		log.Println(err)
		return
	}
	cmd := exec.Command("ffmpeg", "-i", c.st.Path+"/"+fileID+"."+sourceFormat,
		c.st.Path+"/"+targetFileID.String()+"."+targetFormat)
	if err := cmd.Run(); err != nil {
		if err1 := c.repo.UpdateRequest(requestID, status[2], ""); err1 != nil {
			log.Println(err1)
		}
		log.Println(err)
		return
	}

	targetID, err := c.repo.InsertAudio(filename, targetFormat, targetFileID.String())
	if err != nil {
		if err1 := c.repo.UpdateRequest(requestID, status[2], ""); err1 != nil {
			log.Println(err1)
		}
		log.Println(err)
		return
	}

	if err := c.repo.UpdateRequest(requestID, status[1], targetID); err != nil {
		if err1 := c.repo.UpdateRequest(requestID, status[2], ""); err1 != nil {
			log.Println(err1)
		}
		log.Println(err)
		return
	}

	c.sem.Release(1)
}
