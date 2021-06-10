// Package converter implements a function for audio converting.
package converter

import (
	"context"

	"fmt"
	"os"
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
	if err := c.sem.Acquire(context.Background(), 1); err != nil {
		if updateErr := c.repo.UpdateRequest(requestID, status[2], ""); updateErr != nil {
			logger.Error(ctx, updateErr)
		}
		logger.Error(ctx, fmt.Errorf("can't aquire semaphore, %v", err))
		return
	}

	if err := c.repo.UpdateRequest(requestID, status[0], ""); err != nil {
		if updateErr := c.repo.UpdateRequest(requestID, status[2], ""); updateErr != nil {
			logger.Error(ctx, updateErr)
		}
		logger.Error(ctx, fmt.Errorf("can't update request, %v", err))
		return
	}
	logger.Debug(ctx, "status changed to processing")

	targetFileID, err := uuid.NewRandom()
	if err != nil {
		if updateErr := c.repo.UpdateRequest(requestID, status[2], ""); updateErr != nil {
			logger.Error(ctx, updateErr)
		}
		logger.Error(ctx, fmt.Errorf("can't generate target file uuid, %v", err))
		return
	}
	targetFileIDStr := targetFileID.String()

	sourceLocation := fmt.Sprintf(storage.LocationTemplate, fileID, sourceFormat)
	targetLocation := fmt.Sprintf(storage.LocationTemplate, targetFileIDStr, targetFormat)

	cmd := exec.Command("ffmpeg", "-i", sourceLocation, targetLocation)
	if err = cmd.Run(); err != nil {
		if updateErr := c.repo.UpdateRequest(requestID, status[2], ""); updateErr != nil {
			logger.Error(ctx, updateErr)
		}
		logger.Error(ctx, fmt.Errorf("can't perform conversion"))
		return
	}
	logger.Debug(ctx, "convertion performed successfully")

	targetFile, err := os.Open(targetLocation)
	if err != nil {
		if updateErr := c.repo.UpdateRequest(requestID, status[2], ""); updateErr != nil {
			logger.Error(ctx, updateErr)
		}
		logger.Error(ctx, fmt.Errorf("can't generate targetFileID, %v", err))
		return
	}

	err = c.storage.UploadFileToCloud(targetFile, targetFileIDStr, targetFormat)
	if err != nil {
		if updateErr := c.repo.UpdateRequest(requestID, status[2], ""); updateErr != nil {
			logger.Error(ctx, updateErr)
		}
		logger.Error(ctx, fmt.Errorf("can't upload file to s3, %v", err))
		return
	}
	logger.Debug(ctx, "converted file uploaded to s3 successfully")

	targetID, err := c.repo.InsertAudio(filename, targetFormat, targetFileIDStr)
	if err != nil {
		if updateErr := c.repo.UpdateRequest(requestID, status[2], ""); updateErr != nil {
			logger.Error(ctx, updateErr)
		}
		logger.Error(ctx, fmt.Errorf("can't insert audio, %v", err))
		return
	}

	if err := c.repo.UpdateRequest(requestID, status[1], targetID); err != nil {
		if updateErr := c.repo.UpdateRequest(requestID, status[2], ""); updateErr != nil {
			logger.Error(ctx, updateErr)
		}
		logger.Error(ctx, fmt.Errorf("can't update request, %v", err))
		return
	}
	logger.Debug(ctx, "status changed to done")

	logger.Info(ctx, fmt.Sprintf("the file with ID=%s was converted successfully", targetID))
	c.sem.Release(1)
}
