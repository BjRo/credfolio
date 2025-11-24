package handler_test

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/credfolio/apps/backend/internal/domain"
	"github.com/credfolio/apps/backend/internal/handler"
	"github.com/credfolio/apps/backend/internal/handler/middleware"
	"github.com/credfolio/apps/backend/internal/repository/mocks"
	"github.com/credfolio/apps/backend/internal/service"
	servicemocks "github.com/credfolio/apps/backend/internal/service/mocks"
	"github.com/credfolio/apps/backend/pkg/logger"
	pdfmocks "github.com/credfolio/apps/backend/pkg/pdf/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUploadReferenceLetter(t *testing.T) {
	mockProfileRepo := new(mocks.MockProfileRepository)
	mockWorkExpRepo := new(mocks.MockWorkExperienceRepository)
	mockCredibilityRepo := new(mocks.MockCredibilityHighlightRepository)
	mockRefLetterRepo := new(mocks.MockReferenceLetterRepository)
	mockJobMatchRepo := new(mocks.MockJobMatchRepository)
	mockLLMProvider := new(servicemocks.MockLLMProvider)
	mockPDFExtractor := new(pdfmocks.MockPDFExtractor)
	appLogger := logger.New()

	svc := service.NewProfileService(
		mockProfileRepo,
		mockWorkExpRepo,
		mockCredibilityRepo,
		mockRefLetterRepo,
		mockLLMProvider,
		mockPDFExtractor,
		appLogger,
	)

	tailoringSvc := service.NewTailoringService(
		mockProfileRepo,
		mockJobMatchRepo,
		mockLLMProvider,
		appLogger,
	)

	api := handler.NewAPI(svc, tailoringSvc, mockRefLetterRepo, mockJobMatchRepo, mockPDFExtractor, appLogger)

	userID := uuid.New()

	t.Run("Upload Success", func(t *testing.T) {
		mockRefLetterRepo.On("Create", mock.Anything, mock.MatchedBy(func(l *domain.ReferenceLetter) bool {
			return l.UserID == userID && l.ExtractedText == "Reference letter content"
		})).Return(nil)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", "test.txt")
		assert.NoError(t, err)
		textContent := "Reference letter content"
		_, err = part.Write([]byte(textContent))
		assert.NoError(t, err)
		err = writer.Close()
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/reference-letters", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		ctx := middleware.WithUserID(req.Context(), userID)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		api.UploadReferenceLetter(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockRefLetterRepo.AssertExpectations(t)
	})

	t.Run("Invalid File Type", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", "test.pdf")
		assert.NoError(t, err)
		_, err = part.Write([]byte("not a valid file"))
		assert.NoError(t, err)
		err = writer.Close()
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/reference-letters", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		ctx := middleware.WithUserID(req.Context(), userID)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		api.UploadReferenceLetter(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var resp handler.ErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, handler.ErrorCodeInvalidFileType, resp.ErrorID)
		assert.Contains(t, resp.Message, "txt or md")
	})

	t.Run("File Too Large", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", "large.txt")
		assert.NoError(t, err)

		largeContent := make([]byte, 11*1024*1024)
		_, err = part.Write(largeContent)
		assert.NoError(t, err)
		err = writer.Close()
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/reference-letters", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		ctx := middleware.WithUserID(req.Context(), userID)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		api.UploadReferenceLetter(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var resp handler.ErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, handler.ErrorCodeFileTooLarge, resp.ErrorID)
		assert.Contains(t, resp.Message, "exceeds maximum")
	})
}

func TestGenerateProfile(t *testing.T) {
	mockProfileRepo := new(mocks.MockProfileRepository)
	mockWorkExpRepo := new(mocks.MockWorkExperienceRepository)
	mockCredibilityRepo := new(mocks.MockCredibilityHighlightRepository)
	mockRefLetterRepo := new(mocks.MockReferenceLetterRepository)
	mockJobMatchRepo := new(mocks.MockJobMatchRepository)
	mockLLMProvider := new(servicemocks.MockLLMProvider)
	mockPDFExtractor := new(pdfmocks.MockPDFExtractor)
	appLogger := logger.New()

	svc := service.NewProfileService(
		mockProfileRepo,
		mockWorkExpRepo,
		mockCredibilityRepo,
		mockRefLetterRepo,
		mockLLMProvider,
		mockPDFExtractor,
		appLogger,
	)

	tailoringSvc := service.NewTailoringService(
		mockProfileRepo,
		mockJobMatchRepo,
		mockLLMProvider,
		appLogger,
	)

	api := handler.NewAPI(svc, tailoringSvc, mockRefLetterRepo, mockJobMatchRepo, mockPDFExtractor, appLogger)

	userID := uuid.New()
	letterID := uuid.New()
	profileID := uuid.New()

	t.Run("Generate Success", func(t *testing.T) {
		mockRefLetterRepo.On("GetByUserID", mock.Anything, userID).Return([]*domain.ReferenceLetter{
			{ID: letterID, UserID: userID, ExtractedText: "Text"},
		}, nil)

		mockProfileRepo.On("GetByUserID", mock.Anything, userID).Return(&domain.Profile{ID: profileID}, nil)
		mockRefLetterRepo.On("GetByID", mock.Anything, letterID).Return(&domain.ReferenceLetter{ID: letterID, ExtractedText: "Text"}, nil)
		mockLLMProvider.On("ExtractProfileData", mock.Anything, "Text").Return(&service.ExtractedProfileData{
			CompanyName: "A",
			StartDate:   "2022-01-01",
		}, nil)
		mockLLMProvider.On("ExtractCredibility", mock.Anything, "Text").Return(&service.CredibilityData{}, nil)
		mockWorkExpRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
		mockProfileRepo.On("GetByID", mock.Anything, profileID).Return(&domain.Profile{
			ID:      profileID,
			Summary: "Test Summary",
		}, nil)

		req := httptest.NewRequest("POST", "/profile/generate", nil)
		ctx := middleware.WithUserID(req.Context(), userID)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		api.GenerateProfile(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "Test Summary", resp["summary"])
	})
}
