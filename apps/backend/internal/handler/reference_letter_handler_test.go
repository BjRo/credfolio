package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/credfolio/apps/backend/internal/domain"
	"github.com/credfolio/apps/backend/internal/handler"
	"github.com/credfolio/apps/backend/internal/handler/middleware"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockReferenceLetterRepository is a mock implementation
type MockReferenceLetterRepository struct {
	mock.Mock
}

func (m *MockReferenceLetterRepository) Create(ctx context.Context, letter *domain.ReferenceLetter) error {
	args := m.Called(ctx, letter)
	return args.Error(0)
}

func (m *MockReferenceLetterRepository) Update(ctx context.Context, letter *domain.ReferenceLetter) error {
	args := m.Called(ctx, letter)
	return args.Error(0)
}

func (m *MockReferenceLetterRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.ReferenceLetter, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ReferenceLetter), args.Error(1)
}

func (m *MockReferenceLetterRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.ReferenceLetter, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.ReferenceLetter), args.Error(1)
}

func (m *MockReferenceLetterRepository) FindPendingByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.ReferenceLetter, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.ReferenceLetter), args.Error(1)
}

func (m *MockReferenceLetterRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func createMultipartRequest(t *testing.T, fieldName, filename, content string) (*http.Request, *bytes.Buffer) {
	t.Helper()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(fieldName, filename)
	require.NoError(t, err)

	_, err = io.WriteString(part, content)
	require.NoError(t, err)

	err = writer.Close()
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/reference-letters", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req, body
}

// T053: Unit test for ReferenceLetterHandler when uploading valid file returns reference letter ID
func TestReferenceLetterHandler_Upload_WhenValidFile_ReturnsReferenceLetterID(t *testing.T) {
	// Arrange
	mockRepo := new(MockReferenceLetterRepository)
	tempDir := t.TempDir()

	h := handler.NewReferenceLetterHandler(mockRepo, tempDir)

	userID := uuid.New()
	req, _ := createMultipartRequest(t, "file", "test-letter.txt", "This is a reference letter content.")

	// Add user ID to context
	ctx := middleware.SetUserID(req.Context(), userID)
	req = req.WithContext(ctx)

	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.ReferenceLetter")).Return(nil)

	recorder := httptest.NewRecorder()

	// Act
	h.Upload(recorder, req)

	// Assert
	assert.Equal(t, http.StatusCreated, recorder.Code)

	var response handler.ReferenceLetterResponse
	err := json.NewDecoder(recorder.Body).Decode(&response)
	require.NoError(t, err)

	assert.NotEmpty(t, response.ID)
	assert.Equal(t, "test-letter.txt", response.FileName)
	assert.Equal(t, "PENDING", response.Status)
}

// T054: Unit test for ReferenceLetterHandler when uploading invalid file type returns error
func TestReferenceLetterHandler_Upload_WhenInvalidFileType_ReturnsError(t *testing.T) {
	// Arrange
	mockRepo := new(MockReferenceLetterRepository)
	tempDir := t.TempDir()

	h := handler.NewReferenceLetterHandler(mockRepo, tempDir)

	userID := uuid.New()
	req, _ := createMultipartRequest(t, "file", "test-letter.pdf", "PDF content")

	// Add user ID to context
	ctx := middleware.SetUserID(req.Context(), userID)
	req = req.WithContext(ctx)

	recorder := httptest.NewRecorder()

	// Act
	h.Upload(recorder, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "Unsupported file type")
}

func TestReferenceLetterHandler_Upload_WhenNoUserID_ReturnsUnauthorized(t *testing.T) {
	// Arrange
	mockRepo := new(MockReferenceLetterRepository)
	tempDir := t.TempDir()

	h := handler.NewReferenceLetterHandler(mockRepo, tempDir)

	req, _ := createMultipartRequest(t, "file", "test-letter.txt", "Content")

	recorder := httptest.NewRecorder()

	// Act
	h.Upload(recorder, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
}

func TestReferenceLetterHandler_Upload_WhenNoFile_ReturnsBadRequest(t *testing.T) {
	// Arrange
	mockRepo := new(MockReferenceLetterRepository)
	tempDir := t.TempDir()

	h := handler.NewReferenceLetterHandler(mockRepo, tempDir)

	userID := uuid.New()
	req := httptest.NewRequest(http.MethodPost, "/reference-letters", nil)
	ctx := middleware.SetUserID(req.Context(), userID)
	req = req.WithContext(ctx)

	recorder := httptest.NewRecorder()

	// Act
	h.Upload(recorder, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestReferenceLetterHandler_List_WhenLettersExist_ReturnsLetters(t *testing.T) {
	// Arrange
	mockRepo := new(MockReferenceLetterRepository)
	tempDir := t.TempDir()

	h := handler.NewReferenceLetterHandler(mockRepo, tempDir)

	userID := uuid.New()
	letters := []*domain.ReferenceLetter{
		domain.NewReferenceLetter(userID, "letter1.txt", "/path/letter1.txt"),
		domain.NewReferenceLetter(userID, "letter2.txt", "/path/letter2.txt"),
	}

	mockRepo.On("FindByUserID", mock.Anything, userID).Return(letters, nil)

	req := httptest.NewRequest(http.MethodGet, "/reference-letters", nil)
	ctx := middleware.SetUserID(req.Context(), userID)
	req = req.WithContext(ctx)

	recorder := httptest.NewRecorder()

	// Act
	h.List(recorder, req)

	// Assert
	assert.Equal(t, http.StatusOK, recorder.Code)

	var response []handler.ReferenceLetterResponse
	err := json.NewDecoder(recorder.Body).Decode(&response)
	require.NoError(t, err)

	assert.Len(t, response, 2)
}

func TestReferenceLetterHandler_Upload_WhenMarkdownFile_ReturnsSuccess(t *testing.T) {
	// Arrange
	mockRepo := new(MockReferenceLetterRepository)
	tempDir := t.TempDir()

	h := handler.NewReferenceLetterHandler(mockRepo, tempDir)

	userID := uuid.New()
	req, _ := createMultipartRequest(t, "file", "test-letter.md", "# Reference Letter\n\nThis is a test.")

	ctx := middleware.SetUserID(req.Context(), userID)
	req = req.WithContext(ctx)

	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.ReferenceLetter")).Return(nil)

	recorder := httptest.NewRecorder()

	// Act
	h.Upload(recorder, req)

	// Assert
	assert.Equal(t, http.StatusCreated, recorder.Code)
}

func TestReferenceLetterHandler_Upload_WhenRepoFails_ReturnsInternalServerError(t *testing.T) {
	// Arrange
	mockRepo := new(MockReferenceLetterRepository)
	tempDir := t.TempDir()

	h := handler.NewReferenceLetterHandler(mockRepo, tempDir)

	userID := uuid.New()
	req, _ := createMultipartRequest(t, "file", "test-letter.txt", "Content")

	ctx := middleware.SetUserID(req.Context(), userID)
	req = req.WithContext(ctx)

	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.ReferenceLetter")).Return(assert.AnError)

	recorder := httptest.NewRecorder()

	// Act
	h.Upload(recorder, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}
