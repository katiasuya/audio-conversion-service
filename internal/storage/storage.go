// Package storage provides logic to communicate with aws s3 cloud object storage.
package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
)

// Storage represents aws s3 client configuration.
type Storage struct {
	bucket string
	sess   *session.Session
}

// New creates a new storage with the given configuration.
func New(bucket string, sess *session.Session) *Storage {
	return &Storage{
		bucket: bucket,
		sess:   sess,
	}
}

// GetClientConfig gets client configuration.
func (s *Storage) GetClientConfig() (*session.Session, string) {
	return s.sess, s.bucket
}

// UploadFile stores request file in the storage.
func (s *Storage) UploadFile(sourceFile io.Reader, format string) (string, error) {
	fileID, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("can't generate file uuid, %w", err)
	}

	uploader := s3manager.NewUploader(s.sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.bucket),
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
	svc := s3.New(s.sess)

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(fileID + "." + format),
	})
	urlStr, err := req.Presign(15 * time.Minute)
	if err != nil {
		return "", fmt.Errorf("can't create requets's presigned URL, %w", err)
	}

	return urlStr, err

}
