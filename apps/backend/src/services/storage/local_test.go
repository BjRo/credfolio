package storage

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestLocalStorage(t *testing.T) {
	tempDir := t.TempDir()
	storage, err := NewLocalStorage(tempDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}

	filename := "test.txt"
	content := []byte("hello world")

	path, err := storage.Save(filename, bytes.NewReader(content))
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("file not created")
	}

	// Verify content
	readContent, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	if string(readContent) != string(content) {
		t.Errorf("content mismatch: got %s, want %s", readContent, content)
	}

	// Test Get
	r, err := storage.Get(filename)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if c, ok := r.(io.Closer); ok {
		defer func() { _ = c.Close() }()
	}
}
