package storage

import (
	"io"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	BaseDir string
}

func NewLocalStorage(baseDir string) (*LocalStorage, error) {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, err
	}
	return &LocalStorage{BaseDir: baseDir}, nil
}

func (s *LocalStorage) Save(filename string, content io.Reader) (string, error) {
	// Prevent directory traversal
	cleanName := filepath.Base(filename)
	path := filepath.Join(s.BaseDir, cleanName)

	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := io.Copy(f, content); err != nil {
		return "", err
	}

	return path, nil
}

func (s *LocalStorage) Get(filename string) (io.Reader, error) {
	cleanName := filepath.Base(filename)
	path := filepath.Join(s.BaseDir, cleanName)
	return os.Open(path)
}

