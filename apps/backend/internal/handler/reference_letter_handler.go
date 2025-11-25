package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/credfolio/apps/backend/internal/domain"
	"github.com/credfolio/apps/backend/internal/handler/middleware"
	"github.com/credfolio/apps/backend/internal/repository"
	"github.com/credfolio/apps/backend/pkg/extraction"
	"github.com/credfolio/apps/backend/pkg/logger"
	"github.com/google/uuid"
)

const (
	// MaxFileSize is the maximum allowed file size (5MB)
	MaxFileSize = 5 << 20
	// MaxFilenameLength is the maximum allowed filename length
	MaxFilenameLength = 255
)

// AllowedFileExtensions defines the allowed file extensions for reference letters
var AllowedFileExtensions = map[string]bool{
	".txt":      true,
	".md":       true,
	".markdown": true,
}

// AllowedMIMETypes defines the allowed MIME types for reference letters
var AllowedMIMETypes = map[string]bool{
	"text/plain":               true,
	"text/markdown":            true,
	"text/x-markdown":          true,
	"application/octet-stream": true, // Fallback for unknown types
}

// ReferenceLetterHandler handles reference letter HTTP requests
type ReferenceLetterHandler struct {
	repo        repository.ReferenceLetterRepository
	extractor   *extraction.Extractor
	uploadPath  string
	logger      *logger.Logger
	maxFileSize int64
}

// NewReferenceLetterHandler creates a new reference letter handler
func NewReferenceLetterHandler(repo repository.ReferenceLetterRepository, uploadPath string) *ReferenceLetterHandler {
	// Ensure upload directory exists
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		logger.Error("Failed to create upload directory: %v", err)
	}

	return &ReferenceLetterHandler{
		repo:        repo,
		extractor:   extraction.NewExtractor(),
		uploadPath:  uploadPath,
		logger:      logger.New(logger.LevelInfo),
		maxFileSize: MaxFileSize,
	}
}

// ReferenceLetterResponse represents the API response for a reference letter
type ReferenceLetterResponse struct {
	ID         string `json:"id"`
	FileName   string `json:"fileName"`
	UploadDate string `json:"uploadDate"`
	Status     string `json:"status"`
}

// Upload handles POST /reference-letters - file upload
func (h *ReferenceLetterHandler) Upload(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		h.logger.Error("User ID not found in context")
		h.writeError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	h.logger.Info("Processing reference letter upload for user %s", userID.String())

	// Parse multipart form with size limit
	r.Body = http.MaxBytesReader(w, r.Body, h.maxFileSize+1024) // Allow some overhead for form fields
	if err := r.ParseMultipartForm(h.maxFileSize); err != nil {
		h.logger.Error("Failed to parse multipart form: %v", err)
		if err.Error() == "http: request body too large" {
			h.writeError(w, fmt.Sprintf("File too large. Maximum size is %dMB", h.maxFileSize/(1<<20)), http.StatusRequestEntityTooLarge)
			return
		}
		h.writeError(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// Get file from form
	file, header, err := r.FormFile("file")
	if err != nil {
		h.logger.Error("Failed to get file from form: %v", err)
		h.writeError(w, "File is required. Please upload a reference letter file.", http.StatusBadRequest)
		return
	}
	defer func() { _ = file.Close() }()

	// Validate filename
	if err := h.validateFilename(header.Filename); err != nil {
		h.logger.Error("Invalid filename: %s - %v", header.Filename, err)
		h.writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate file size
	if header.Size > h.maxFileSize {
		h.logger.Error("File too large: %d bytes (max %d)", header.Size, h.maxFileSize)
		h.writeError(w, fmt.Sprintf("File too large. Maximum size is %dMB", h.maxFileSize/(1<<20)), http.StatusRequestEntityTooLarge)
		return
	}

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !AllowedFileExtensions[ext] {
		h.logger.Error("Unsupported file type: %s", ext)
		h.writeError(w, "Unsupported file type. Allowed formats: .txt, .md, .markdown", http.StatusBadRequest)
		return
	}

	// Validate MIME type (content sniffing)
	contentType := header.Header.Get("Content-Type")
	if contentType != "" && !AllowedMIMETypes[contentType] {
		h.logger.Warn("Unexpected MIME type: %s for file %s, proceeding with extension validation", contentType, header.Filename)
	}

	// Read file content to validate it's actually text
	content, err := io.ReadAll(io.LimitReader(file, h.maxFileSize))
	if err != nil {
		h.logger.Error("Failed to read file: %v", err)
		h.writeError(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	// Validate content is valid UTF-8 text
	if !isValidUTF8Text(content) {
		h.logger.Error("File content is not valid UTF-8 text")
		h.writeError(w, "File must contain valid text content", http.StatusBadRequest)
		return
	}

	// Check for empty files
	if len(strings.TrimSpace(string(content))) == 0 {
		h.logger.Error("File is empty or contains only whitespace")
		h.writeError(w, "File is empty. Please upload a reference letter with content.", http.StatusBadRequest)
		return
	}

	// Generate unique filename (sanitized)
	fileID := uuid.New()
	safeFilename := h.sanitizeFilename(header.Filename)
	filename := fileID.String() + "-" + safeFilename
	storagePath := filepath.Join(h.uploadPath, filename)

	// Write content to destination file
	if err := os.WriteFile(storagePath, content, 0644); err != nil {
		h.logger.Error("Failed to write file: %v", err)
		h.writeError(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	// Extract text from file
	extractedText, err := h.extractor.ExtractFromFile(storagePath)
	if err != nil {
		h.logger.Error("Failed to extract text: %v", err)
		// Continue - we can try again later
	}

	// Create reference letter record
	letter := domain.NewReferenceLetter(userID, header.Filename, storagePath)
	if extractedText != "" {
		letter.ExtractedText = extractedText
	}

	if err := h.repo.Create(r.Context(), letter); err != nil {
		h.logger.Error("Failed to create reference letter: %v", err)
		// Clean up uploaded file on database error
		_ = os.Remove(storagePath)
		h.writeError(w, "Failed to save reference letter", http.StatusInternalServerError)
		return
	}

	h.logger.Info("Reference letter uploaded successfully: %s for user %s (size: %d bytes)", letter.ID.String(), userID.String(), len(content))

	// Return response
	response := ReferenceLetterResponse{
		ID:         letter.ID.String(),
		FileName:   letter.FileName,
		UploadDate: letter.UploadDate.Format("2006-01-02T15:04:05Z07:00"),
		Status:     string(letter.Status),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response)
}

// validateFilename checks if the filename is valid and safe
func (h *ReferenceLetterHandler) validateFilename(filename string) error {
	if filename == "" {
		return fmt.Errorf("filename is required")
	}
	if len(filename) > MaxFilenameLength {
		return fmt.Errorf("filename too long (max %d characters)", MaxFilenameLength)
	}
	// Check for path traversal attempts
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		return fmt.Errorf("invalid filename")
	}
	return nil
}

// sanitizeFilename removes potentially dangerous characters from filename
func (h *ReferenceLetterHandler) sanitizeFilename(filename string) string {
	// Remove path separators and null bytes
	filename = strings.ReplaceAll(filename, "/", "_")
	filename = strings.ReplaceAll(filename, "\\", "_")
	filename = strings.ReplaceAll(filename, "\x00", "")
	// Keep only alphanumeric, dots, dashes, and underscores
	var result strings.Builder
	for _, r := range filename {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '.' || r == '-' || r == '_' {
			result.WriteRune(r)
		}
	}
	sanitized := result.String()
	if sanitized == "" {
		return "upload.txt"
	}
	return sanitized
}

// isValidUTF8Text checks if the content is valid UTF-8 text (not binary)
func isValidUTF8Text(content []byte) bool {
	// Check for null bytes (binary indicator)
	for _, b := range content {
		if b == 0 {
			return false
		}
	}
	// Check if content is valid UTF-8
	return len(content) == 0 || strings.ToValidUTF8(string(content), "") == string(content)
}

// writeError writes a JSON error response
func (h *ReferenceLetterHandler) writeError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// List handles GET /reference-letters - list user's reference letters
func (h *ReferenceLetterHandler) List(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		h.logger.Error("User ID not found in context for list request")
		h.writeError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	h.logger.Info("Listing reference letters for user %s", userID.String())

	letters, err := h.repo.FindByUserID(r.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to find reference letters for user %s: %v", userID.String(), err)
		h.writeError(w, "Failed to get reference letters", http.StatusInternalServerError)
		return
	}

	var response []ReferenceLetterResponse
	for _, letter := range letters {
		response = append(response, ReferenceLetterResponse{
			ID:         letter.ID.String(),
			FileName:   letter.FileName,
			UploadDate: letter.UploadDate.Format("2006-01-02T15:04:05Z07:00"),
			Status:     string(letter.Status),
		})
	}

	h.logger.Info("Found %d reference letters for user %s", len(response), userID.String())

	w.Header().Set("Content-Type", "application/json")
	if response == nil {
		response = []ReferenceLetterResponse{} // Return empty array instead of null
	}
	_ = json.NewEncoder(w).Encode(response)
}
