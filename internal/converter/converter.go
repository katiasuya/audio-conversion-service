// Package converter implements a function for audio converting.
package converter

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
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
func (c *Converter) Convert(fileID, filename, sourceFormat, targetFormat, requestID string) {
	if err := c.sem.Acquire(context.Background(), 1); err != nil {
		if err1 := c.repo.UpdateRequest(requestID, status[2], ""); err1 != nil {
			log.Println(err1)
		}
		log.Println(fmt.Errorf("can't aquire semaphore, %w", err))
		return
	}

	if err := c.repo.UpdateRequest(requestID, status[0], ""); err != nil {
		if err1 := c.repo.UpdateRequest(requestID, status[2], ""); err1 != nil {
			log.Println(err1)
		}
		log.Println(fmt.Errorf("can't update request, %w", err))
		return
	}

	targetFileID, err := uuid.NewRandom()
	if err != nil {
		if err1 := c.repo.UpdateRequest(requestID, status[2], ""); err1 != nil {
			log.Println(err1)
		}
		log.Println(fmt.Errorf("can't generate target file uuid, %w", err))
		return
	}

	cmd := exec.Command("ffmpeg", "-i", "/tmp/"+fileID+"."+sourceFormat, "/tmp/"+targetFileID.String()+"."+targetFormat)
	if err := cmd.Run(); err != nil {
		if err1 := c.repo.UpdateRequest(requestID, status[2], ""); err1 != nil {
			log.Println(err1)
		}
		log.Println(fmt.Errorf("can't perform conversion, %w", err))
		return
	}

	targetFile, err := os.Open("/tmp/" + targetFileID.String() + "." + targetFormat)
	if err != nil {
		if err1 := c.repo.UpdateRequest(requestID, status[2], ""); err1 != nil {
			log.Println(err1)
		}
		log.Println(fmt.Errorf("can't open file, %w", err))
		return
	}
	sess, bucket := c.storage.GetClientConfig()
	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(targetFileID.String() + "." + targetFormat),
		Body:   targetFile,
	})
	if err != nil {
		log.Println(fmt.Errorf("can't upload file, %w", err))
		return
	}

	targetID, err := c.repo.InsertAudio(filename, targetFormat, targetFileID.String())
	if err != nil {
		if err1 := c.repo.UpdateRequest(requestID, status[2], ""); err1 != nil {
			log.Println(err1)
		}
		log.Println(fmt.Errorf("can't insert audio, %w", err))
		return
	}

	if err := c.repo.UpdateRequest(requestID, status[1], targetID); err != nil {
		if err1 := c.repo.UpdateRequest(requestID, status[2], ""); err1 != nil {
			log.Println(err1)
		}
		log.Println(fmt.Errorf("can't update request, %w", err))
		return
	}

	c.sem.Release(1)
}
