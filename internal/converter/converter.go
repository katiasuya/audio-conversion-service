package converter

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/google/uuid"
	"github.com/katiasuya/audio-conversion-service/internal/repository"
	"github.com/katiasuya/audio-conversion-service/internal/storage"
)

var status = []string{"processing", "done", "failed"}

// Converter converts audio files to other formats.
type Converter struct {
	repo    *repository.Repository
	storage *storage.Storage
}

// New creates a new Converter with given fields.
func New(repo *repository.Repository, storage *storage.Storage) *Converter {
	return &Converter{
		repo:    repo,
		storage: storage,
	}
}

// Process implements audio conversion process.
func (c *Converter) Process(fileID, filename, sourceFormat, targetFormat, requestID string) error {
	if err := c.repo.UpdateRequest(requestID, status[0], ""); err != nil {
		return fmt.Errorf("can't update request: %w", err)
	}

	targetFileID, err := uuid.NewRandom()
	if err != nil {
		if updateErr := c.repo.UpdateRequest(requestID, status[2], ""); updateErr != nil {
			return fmt.Errorf("can't update request: %w", err)
		}
		return fmt.Errorf("can't generate target file uuid: %w", err)
	}
	targetFileIDStr := targetFileID.String()

	sourceLocation := fmt.Sprintf(storage.LocationTemplate, fileID, sourceFormat)
	targetLocation := fmt.Sprintf(storage.LocationTemplate, targetFileIDStr, targetFormat)

	cmd := exec.Command("ffmpeg", "-i", sourceLocation, targetLocation)
	if err = cmd.Run(); err != nil {
		if updateErr := c.repo.UpdateRequest(requestID, status[2], ""); updateErr != nil {
			return fmt.Errorf("can't update request: %w", err)
		}
		return fmt.Errorf("can't perform conversion")
	}

	targetFile, err := os.Open(targetLocation)
	if err != nil {
		if updateErr := c.repo.UpdateRequest(requestID, status[2], ""); updateErr != nil {
			return fmt.Errorf("can't update request: %w", err)
		}
		return fmt.Errorf("can't generate targetFileID: %w", err)
	}

	err = c.storage.UploadFileToCloud(targetFile, targetFileIDStr, targetFormat)
	if err != nil {
		if updateErr := c.repo.UpdateRequest(requestID, status[2], ""); updateErr != nil {
			return fmt.Errorf("can't update request: %w", err)
		}
		return fmt.Errorf("can't upload file to s3: %w", err)
	}

	targetID, err := c.repo.InsertAudio(filename, targetFormat, targetFileIDStr)
	if err != nil {
		if updateErr := c.repo.UpdateRequest(requestID, status[2], ""); updateErr != nil {
			return fmt.Errorf("can't update request: %w", err)
		}
		return fmt.Errorf("can't insert audio: %w", err)
	}

	if err := c.repo.UpdateRequest(requestID, status[1], targetID); err != nil {
		return fmt.Errorf("can't update request: %w", err)
	}

	return nil
}
