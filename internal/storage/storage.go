// Package storage provides logic to communicate with aws s3 cloud object storage.
package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
)

// Storage represents aws s3 configuration.
type Storage struct {
	accessKeyID     string
	secretAccessKey string
	region          string
	Bucket          string
}

// New creates a new storage with the given configuration.
func New(accessKeyID, secretAccessKey, region, bucket string) *Storage {
	return &Storage{
		accessKeyID:     accessKeyID,
		secretAccessKey: secretAccessKey,
		region:          region,
		Bucket:          bucket,
	}
}

// NewSession creates new aws session.
func (s *Storage) NewSession() (*session.Session, error) {
	return session.NewSession(
		&aws.Config{
			Region:      aws.String(s.region),
			Credentials: credentials.NewStaticCredentials(s.accessKeyID, s.secretAccessKey, ""),
		},
	)
}

// UploadFile stores request file in the storage.
func (s *Storage) UploadFile(sourceFile io.Reader, format string) (string, error) {
	sess, err := s.NewSession()
	if err != nil {
		return "", fmt.Errorf("can't create session, %w", err)
	}

	fileID, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("can't generate file uuid, %w", err)
	}

	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(fileID.String() + "." + format),
		Body:   sourceFile,
	})
	if err != nil {
		return "", fmt.Errorf("can't upload file, %w", err)
	}

	file, err := os.Create(filepath.Join("/tmp", fileID.String()+"."+format))
	if err != nil {
		return "", fmt.Errorf("can't create file, %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, sourceFile); err != nil {
		return "", fmt.Errorf("can't copy file, %w", err)
	}

	return fileID.String(), nil
}

// GetDownloadURL generates URL to download the file from the storage.
func (s *Storage) GetDownloadURL(fileID, format string) (string, error) {
	sess, err := s.NewSession()
	if err != nil {
		return "", fmt.Errorf("can't create session, %w", err)
	}
	svc := s3.New(sess)

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(fileID + "." + format),
	})
	urlStr, err := req.Presign(15 * time.Minute)
	if err != nil {
		return "", fmt.Errorf("can't create requets's presigned URL, %w", err)
	}

	return urlStr, err
}
