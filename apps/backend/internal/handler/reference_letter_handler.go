package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/credfolio/apps/backend/api/generated"
	"github.com/credfolio/apps/backend/internal/domain"
	"github.com/credfolio/apps/backend/internal/handler/middleware"
	"github.com/google/uuid"
)

// UploadReferenceLetter implements generated.ServerInterface
func (a *API) UploadReferenceLetter(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == uuid.Nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse multipart form (10MB max)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		a.Logger.Error("Failed to parse multipart form: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		a.Logger.Error("Failed to get file from form: %v", err)
		http.Error(w, "File is required", http.StatusBadRequest)
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			a.Logger.Error("Failed to close uploaded file: %v", err)
		}
	}()

	// Ensure uploads directory exists
	uploadDir := "tmp/uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		a.Logger.Error("Failed to create upload directory: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Save file to disk
	filename := fmt.Sprintf("%s-%s", uuid.New().String(), header.Filename)
	filePath := filepath.Join(uploadDir, filename)

	dst, err := os.Create(filePath)
	if err != nil {
		a.Logger.Error("Failed to create destination file: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := dst.Close(); err != nil {
			a.Logger.Error("Failed to close destination file: %v", err)
		}
	}()

	if _, err := io.Copy(dst, file); err != nil {
		a.Logger.Error("Failed to copy file content: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Extract text from PDF
	// We need to re-open the file or seek to start because io.Copy consumed it
	// But Extractor.ExtractText takes an io.Reader.
	// Since we saved it to disk, we can open the saved file.
	savedFile, err := os.Open(filePath)
	if err != nil {
		a.Logger.Error("Failed to open saved file for extraction: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := savedFile.Close(); err != nil {
			a.Logger.Error("Failed to close saved file: %v", err)
		}
	}()

	extractedText, err := a.PDFExtractor.ExtractText(savedFile)
	if err != nil {
		a.Logger.Error("Failed to extract text from PDF: %v", err)
		http.Error(w, "Failed to process PDF", http.StatusInternalServerError)
		return
	}

	// Create ReferenceLetter record
	status := domain.ReferenceLetterStatusPending
	if extractedText != "" {
		status = domain.ReferenceLetterStatusProcessed
	}

	letter := &domain.ReferenceLetter{
		UserID:        userID,
		FileName:      header.Filename,
		StoragePath:   filePath,
		Status:        status,
		ExtractedText: extractedText,
		UploadDate:    time.Now(),
	}

	if err := a.ReferenceLetterRepo.Create(r.Context(), letter); err != nil {
		a.Logger.Error("Failed to create reference letter record: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	statusStr := string(letter.Status)
	resp := generated.ReferenceLetter{
		Id:         &letter.ID,
		FileName:   &letter.FileName,
		UploadDate: &letter.UploadDate,
		Status:     &statusStr,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		a.Logger.Error("Failed to encode response: %v", err)
	}
}
