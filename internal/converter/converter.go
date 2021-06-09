// Package converter implements a function for audio converting.
package converter

import (
	"context"
	"log"

	"fmt"
	"os"
	"os/exec"

	"github.com/google/uuid"
	"github.com/katiasuya/audio-conversion-service/internal/mycontext"
	"github.com/katiasuya/audio-conversion-service/internal/repository"
	"github.com/katiasuya/audio-conversion-service/internal/storage"
	"golang.org/x/sync/semaphore"
)

var status = []string{"processing", "done", "failed"}

// Converter converts audio files to other formats.
type Converter struct {
	sem     *semaphore.Weighted
	repo    *repository.Repository
	storage *storage.Storage
}

// New creates a new Converter with given fields.
func New(sem *semaphore.Weighted, repo *repository.Repository, storage *storage.Storage) *Converter {
	return &Converter{
		sem:     sem,
		repo:    repo,
		storage: storage,
	}
}

// Convert implements audio conversion.
func (c *Converter) Convert(ctx context.Context, fileID, filename, sourceFormat, targetFormat, requestID string) {
	ctxLogger, ok := mycontext.LoggerFromContext(ctx)
	if !ok {
		log.Fatal("can't get logger from context in Converter")
		return
	}
	logger := ctxLogger.WithField("package", "converter")

	if err := c.sem.Acquire(context.Background(), 1); err != nil {
		if updateErr := c.repo.UpdateRequest(requestID, status[2], ""); updateErr != nil {
			logger.Errorln(updateErr)
		}
		logger.Errorln(fmt.Errorf("can't aquire semaphore, %w", err))
		return
	}

	if err := c.repo.UpdateRequest(requestID, status[0], ""); err != nil {
		if updateErr := c.repo.UpdateRequest(requestID, status[2], ""); updateErr != nil {
			logger.Errorln(updateErr)
		}
		logger.Errorln(fmt.Errorf("can't update request, %w", err))
		return
	}
	logger.Debugln("status changed to processing")

	targetFileID, err := uuid.NewRandom()
	if err != nil {
		if updateErr := c.repo.UpdateRequest(requestID, status[2], ""); updateErr != nil {
			logger.Errorln(updateErr)
		}
		logger.Errorln(fmt.Errorf("can't generate target file uuid, %w", err))
		return
	}
	targetFileIDStr := targetFileID.String()

	sourceLocation := fmt.Sprintf(storage.LocationTemplate, fileID, sourceFormat)
	targetLocation := fmt.Sprintf(storage.LocationTemplate, targetFileIDStr, targetFormat)

	cmd := exec.Command("ffmpeg", "-i", sourceLocation, targetLocation)
	if err = cmd.Run(); err != nil {
		if updateErr := c.repo.UpdateRequest(requestID, status[2], ""); updateErr != nil {
			logger.Errorln(updateErr)
		}
		logger.Errorln(fmt.Errorf("can't perform conversion, %w", err))
		return
	}
	logger.Debugln("convertion performed successfully")

	targetFile, err := os.Open(targetLocation)
	if err != nil {
		if updateErr := c.repo.UpdateRequest(requestID, status[2], ""); updateErr != nil {
			logger.Errorln(updateErr)
		}
		logger.Errorln(fmt.Errorf("can't generate targetFileID, %w", err))
		return
	}

	err = c.storage.UploadFileToCloud(targetFile, targetFileIDStr, targetFormat)
	if err != nil {
		if updateErr := c.repo.UpdateRequest(requestID, status[2], ""); updateErr != nil {
			logger.Errorln(updateErr)
		}
		logger.Errorln("can't upload file to s3: ", err)
		return
	}
	logger.Debugln("converted file uploaded to s3 successfully")

	targetID, err := c.repo.InsertAudio(filename, targetFormat, targetFileIDStr)
	if err != nil {
		if updateErr := c.repo.UpdateRequest(requestID, status[2], ""); updateErr != nil {
			logger.Errorln(updateErr)
		}
		logger.Errorln(fmt.Errorf("can't insert audio, %w", err))
		return
	}

	if err := c.repo.UpdateRequest(requestID, status[1], targetID); err != nil {
		if updateErr := c.repo.UpdateRequest(requestID, status[2], ""); updateErr != nil {
			logger.Errorln(updateErr)
		}
		logger.Errorln("can't update request: ", err)
		return
	}
	logger.Debugln("status changed to done")

	logger.WithField("targetID", targetID).Infoln("the file was converted successfully")
	c.sem.Release(1)
}
