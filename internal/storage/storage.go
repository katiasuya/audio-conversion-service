// Package storage provides logic to communicate with aws s3 cloud object storage.
package storage

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
)

const LocationTemplate = "/tmp/%s.%s"

// Storage represents aws s3 client.
type Storage struct {
	svc      *s3.S3
	bucket   string
	uploader *s3manager.Uploader
}

// New creates a new storage with the given client.
func New(svc *s3.S3, bucket string, uploader *s3manager.Uploader) *Storage {
	return &Storage{
		svc:      svc,
		bucket:   bucket,
		uploader: uploader,
	}
}

// UploadFile uploads request file.
func (s *Storage) UploadFile(sourceFile io.Reader, format string) (string, error) {
	fileID, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("can't generate file uuid, %w", err)
	}
	fileIDStr := fileID.String()

	if err := s.UploadFileToCloud(sourceFile, fileIDStr, format); err != nil {
		return "", err
	}

	file, err := os.Create(fmt.Sprintf(LocationTemplate, fileIDStr, format))
	if err != nil {
		return "", fmt.Errorf("can't create file, %w", err)
	}
	defer file.Close()
	if _, err := io.Copy(file, sourceFile); err != nil {
		return "", fmt.Errorf("can't copy file, %w", err)
	}

	return fileIDStr, nil
}

// UploadFileToCloud uploads request file to s3 cloud storage.
func (s *Storage) UploadFileToCloud(sourceFile io.Reader, fileID, format string) error {
	_, err := s.uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(fileID + "." + format),
		Body:   sourceFile,
	})
	if err != nil {
		return fmt.Errorf("can't upload file to S3, %w", err)
	}
	return nil
}

// GetDownloadURL generates URL to download the file from the storage.
func (s *Storage) GetDownloadURL(fileID, format string) (string, error) {
	req, _ := s.svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(fileID + "." + format),
	})
	urlStr, err := req.Presign(15 * time.Minute)
	if err != nil {
		return "", fmt.Errorf("can't create requets's presigned URL, %w", err)
	}

	return urlStr, err
}
