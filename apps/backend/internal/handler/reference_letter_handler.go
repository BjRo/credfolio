package handler

import (
	"encoding/json"
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

// ReferenceLetterHandler handles reference letter HTTP requests
type ReferenceLetterHandler struct {
	repo       repository.ReferenceLetterRepository
	extractor  *extraction.Extractor
	uploadPath string
	logger     *logger.Logger
}

// NewReferenceLetterHandler creates a new reference letter handler
func NewReferenceLetterHandler(repo repository.ReferenceLetterRepository, uploadPath string) *ReferenceLetterHandler {
	// Ensure upload directory exists
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		logger.Error("Failed to create upload directory: %v", err)
	}

	return &ReferenceLetterHandler{
		repo:       repo,
		extractor:  extraction.NewExtractor(),
		uploadPath: uploadPath,
		logger:     logger.New(logger.LevelInfo),
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
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse multipart form (max 10MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		h.logger.Error("Failed to parse multipart form: %v", err)
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get file from form
	file, header, err := r.FormFile("file")
	if err != nil {
		h.logger.Error("Failed to get file from form: %v", err)
		http.Error(w, "File is required", http.StatusBadRequest)
		return
	}
	defer func() { _ = file.Close() }()

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !h.extractor.IsSupportedExtension(ext) {
		h.logger.Error("Unsupported file type: %s", ext)
		http.Error(w, "Unsupported file type. Supported: .txt, .md, .markdown", http.StatusBadRequest)
		return
	}

	// Generate unique filename
	fileID := uuid.New()
	filename := fileID.String() + "-" + header.Filename
	storagePath := filepath.Join(h.uploadPath, filename)

	// Create destination file
	dst, err := os.Create(storagePath)
	if err != nil {
		h.logger.Error("Failed to create file: %v", err)
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer func() { _ = dst.Close() }()

	// Copy uploaded file to destination
	if _, err := io.Copy(dst, file); err != nil {
		h.logger.Error("Failed to copy file: %v", err)
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
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
		http.Error(w, "Failed to save reference letter", http.StatusInternalServerError)
		return
	}

	h.logger.Info("Reference letter uploaded: %s for user %s", letter.ID.String(), userID.String())

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

// List handles GET /reference-letters - list user's reference letters
func (h *ReferenceLetterHandler) List(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	letters, err := h.repo.FindByUserID(r.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to find reference letters: %v", err)
		http.Error(w, "Failed to get reference letters", http.StatusInternalServerError)
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

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}
