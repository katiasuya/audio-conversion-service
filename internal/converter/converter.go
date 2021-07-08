package converter

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/google/uuid"
	"github.com/katiasuya/audio-conversion-service/internal/repository"
	"github.com/katiasuya/audio-conversion-service/internal/storage"
)

// Conversion statuses.
const (
	statusProcessing = "processing"
	statusDone       = "done"
	statusFailed     = "failed"
)

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
	err := c.convert(fileID, filename, sourceFormat, targetFormat, requestID)
	if err != nil {
		updateErr := c.repo.UpdateRequest(requestID, statusFailed, "")
		if updateErr != nil {
			return fmt.Errorf("can't update request: %w", err)
		}
		return err
	}
	return nil
}

func (c *Converter) convert(fileID, filename, sourceFormat, targetFormat, requestID string) error {
	err := c.repo.UpdateRequest(requestID, statusProcessing, "")
	if err != nil {
		return fmt.Errorf("can't update request: %w", err)
	}

	err = c.storage.DownloadFileFromCloud(fileID, sourceFormat)
	if err != nil {
		return err
	}

	targetFileID := uuid.NewString()

	sourceLocation := fmt.Sprintf(storage.LocationTmpl, fileID, sourceFormat)
	targetLocation := fmt.Sprintf(storage.LocationTmpl, targetFileID, targetFormat)

	cmd := exec.Command("ffmpeg", "-i", sourceLocation, targetLocation)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("can't perform conversion")
	}

	targetFile, err := os.Open(targetLocation)
	if err != nil {
		return fmt.Errorf("can't generate targetFileID: %w", err)
	}

	err = c.storage.UploadFileToCloud(targetFile, targetFileID, targetFormat)
	if err != nil {
		return fmt.Errorf("can't upload file to s3: %w", err)
	}

	targetID, err := c.repo.InsertAudio(filename, targetFormat, targetFileID)
	if err != nil {
		return fmt.Errorf("can't insert audio: %w", err)
	}

	err = c.repo.UpdateRequest(requestID, statusDone, targetID)
	if err != nil {
		return fmt.Errorf("can't update request: %w", err)
	}

	return nil
}
