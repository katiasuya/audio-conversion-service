// Package storage provides logic to communicate with file storage.
package storage

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

// Storage represents a directory path where all the files will be downloaded to.
type Storage struct {
	Path string
}

// New creates a new storage in the given directory path.
func New(path string) *Storage {
	return &Storage{
		Path: path,
	}
}

// UploadFile stores request file in the storage.
func (s *Storage) UploadFile(sourceFile multipart.File, format string) (string, error) {
	fileID, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	fileLocation := filepath.Join(s.Path, fileID.String()+"."+format)

	file, err := os.Create(fileLocation)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if _, err := io.Copy(file, sourceFile); err != nil {
		return "", err
	}

	return fileID.String(), nil
}

// DownloadFile downloads the file by its id from the storage.
func (s *Storage) DownloadFile(fileID, format string) (multipart.File, error) {
	fileLocation := filepath.Join(s.Path, fileID+"."+format)

	_, err := os.Stat(fileLocation)
	if err != nil {
		return nil, err
	}

	return os.Open(fileLocation)
}
