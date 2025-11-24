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
		writeErrorResponse(w, http.StatusUnauthorized, ErrorCodeUnauthorized, "Unauthorized")
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		a.Logger.Error("Failed to parse multipart form: %v", err)
		writeErrorResponse(w, http.StatusBadRequest, ErrorCodeInvalidRequestBody, "Invalid request body")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		a.Logger.Error("Failed to get file from form: %v", err)
		writeErrorResponse(w, http.StatusBadRequest, ErrorCodeMissingField, "File is required")
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			a.Logger.Error("Failed to close uploaded file: %v", err)
		}
	}()

	if err := ValidateFileType(header.Header.Get("Content-Type"), header.Filename); err != nil {
		if valErr, ok := err.(*ValidationError); ok {
			a.Logger.Error("Invalid file type: %s (MIME: %s, Filename: %s)", valErr.Message, header.Header.Get("Content-Type"), header.Filename)
			writeErrorResponse(w, http.StatusBadRequest, valErr.ErrorCode, valErr.Message)
			return
		}
		writeErrorResponse(w, http.StatusBadRequest, ErrorCodeInvalidFileType, "Invalid file type")
		return
	}

	const maxFileSize = 10 << 20 // 10MB in bytes
	if err := ValidateFileSize(header.Size, maxFileSize); err != nil {
		if valErr, ok := err.(*ValidationError); ok {
			a.Logger.Error("File size validation failed: %s (Size: %d bytes)", valErr.Message, header.Size)
			writeErrorResponse(w, http.StatusBadRequest, valErr.ErrorCode, valErr.Message)
			return
		}
		writeErrorResponse(w, http.StatusBadRequest, ErrorCodeFileTooLarge, "File size exceeds maximum allowed")
		return
	}

	sanitizedFilename, err := SanitizeFilename(header.Filename)
	if err != nil {
		if valErr, ok := err.(*ValidationError); ok {
			writeErrorResponse(w, http.StatusBadRequest, valErr.ErrorCode, valErr.Message)
			return
		}
		writeErrorResponse(w, http.StatusBadRequest, ErrorCodeInvalidRequest, "Invalid filename")
		return
	}

	uploadDir := "tmp/uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		a.Logger.Error("Failed to create upload directory: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, ErrorCodeInternalServerError, "Internal server error")
		return
	}

	filename := fmt.Sprintf("%s-%s", uuid.New().String(), sanitizedFilename)
	filePath := filepath.Join(uploadDir, filename)

	dst, err := os.Create(filePath)
	if err != nil {
		a.Logger.Error("Failed to create destination file: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, ErrorCodeInternalServerError, "Internal server error")
		return
	}
	defer func() {
		if err := dst.Close(); err != nil {
			a.Logger.Error("Failed to close destination file: %v", err)
		}
	}()

	if _, err := io.Copy(dst, file); err != nil {
		a.Logger.Error("Failed to copy file content: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, ErrorCodeInternalServerError, "Internal server error")
		return
	}

	savedFile, err := os.Open(filePath)
	if err != nil {
		a.Logger.Error("Failed to open saved file for extraction: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, ErrorCodeInternalServerError, "Internal server error")
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
		writeErrorResponse(w, http.StatusInternalServerError, ErrorCodePDFExtractionFailed, "Failed to process PDF")
		return
	}

	status := domain.ReferenceLetterStatusPending
	if extractedText != "" {
		status = domain.ReferenceLetterStatusProcessed
	}

	letter := &domain.ReferenceLetter{
		UserID:        userID,
		FileName:      sanitizedFilename,
		StoragePath:   filePath,
		Status:        status,
		ExtractedText: extractedText,
		UploadDate:    time.Now(),
	}

	if err := a.ReferenceLetterRepo.Create(r.Context(), letter); err != nil {
		a.Logger.Error("Failed to create reference letter record: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, ErrorCodeDatabaseError, "Internal server error")
		return
	}

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
