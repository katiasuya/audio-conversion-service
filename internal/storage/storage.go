// Package storage provides logic to communicate with aws s3 cloud object storage.
package storage

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
	"github.com/katiasuya/audio-conversion-service/internal/config"
)

const (
	filenameTmpl = "%s.%s"
	LocationTmpl = "/tmp/" + filenameTmpl
)

// Storage represents aws s3 client.
type Storage struct {
	svc        *s3.S3
	bucket     string
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
}

// NewS3Client creates new S3 client.
func NewS3Client(conf *config.AWSData) (*Storage, error) {
	sess, err := session.NewSession(
		&aws.Config{
			Region:      aws.String(conf.Region),
			Credentials: credentials.NewStaticCredentials(conf.AccessKeyID, conf.SecretAccessKey, ""),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("can't create new session: %w", err)
	}

	downloader := s3manager.NewDownloader(sess)
	uploader := s3manager.NewUploader(sess)
	svc := s3.New(sess)

	return &Storage{
		svc:        svc,
		bucket:     conf.Bucket,
		uploader:   uploader,
		downloader: downloader,
	}, nil
}

// UploadFile uploads request file.
func (s *Storage) UploadFile(sourceFile io.Reader, format string) (string, error) {
	fileID := uuid.NewString()

	err := s.UploadFileToCloud(sourceFile, fileID, format)
	if err != nil {
		return "", err
	}

	file, err := os.Create(fmt.Sprintf(LocationTmpl, fileID, format))
	if err != nil {
		return "", fmt.Errorf("can't create local file, %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, sourceFile)
	if err != nil {
		return "", fmt.Errorf("can't copy file, %w", err)
	}

	return fileID, nil
}

// UploadFileToCloud uploads request file to s3 cloud storage.
func (s *Storage) UploadFileToCloud(sourceFile io.Reader, fileID, format string) error {
	_, err := s.uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(fmt.Sprintf(filenameTmpl, fileID, format)),
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
		Key:    aws.String(fmt.Sprintf(filenameTmpl, fileID, format)),
	})
	urlStr, err := req.Presign(15 * time.Minute)
	if err != nil {
		return "", fmt.Errorf("can't create requets's presigned URL, %w", err)
	}

	return urlStr, err
}

// DownloadFileFromCloud downloads request file from s3 cloud storage.
func (s *Storage) DownloadFileFromCloud(fileID, format string) error {
	filename := fmt.Sprintf(LocationTmpl, fileID, format)

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("can't create local file, %w", err)
	}
	defer file.Close()

	_, err = s.downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    aws.String(fmt.Sprintf(filenameTmpl, fileID, format)),
		})
	if err != nil {
		return fmt.Errorf("can't download file from S3, %w", err)
	}

	return nil
}
