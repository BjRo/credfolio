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
			return l.UserID == userID
		})).Return(nil)

		mockPDFExtractor.On("ExtractText", mock.Anything).Return("Extracted Text", nil)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", "test.pdf")
		assert.NoError(t, err)
		pdfContent := "%PDF-1.4\n1 0 obj\n<<\n/Type /Catalog\n/Pages 2 0 R\n>>\nendobj\n2 0 obj\n<<\n/Type /Pages\n/Kids [3 0 R]\n/Count 1\n>>\nendobj\n3 0 obj\n<<\n/Type /Page\n/Parent 2 0 R\n/Resources <<\n/Font <<\n/F1 4 0 R\n>>\n>>\n/MediaBox [0 0 612 792]\n/Contents 5 0 R\n>>\nendobj\n4 0 obj\n<<\n/Type /Font\n/Subtype /Type1\n/BaseFont /Helvetica\n>>\nendobj\n5 0 obj\n<< /Length 44 >>\nstream\nBT\n/F1 24 Tf\n100 100 Td\n(Hello World) Tj\nET\nendstream\nendobj\nxref\n0 6\n0000000000 65535 f \n0000000010 00000 n \n0000000060 00000 n \n0000000117 00000 n \n0000000238 00000 n \n0000000325 00000 n \ntrailer\n<<\n/Size 6\n/Root 1 0 R\n>>\nstartxref\n419\n%%EOF"
		_, err = part.Write([]byte(pdfContent))
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
		part, err := writer.CreateFormFile("file", "test.txt")
		assert.NoError(t, err)
		_, err = part.Write([]byte("not a pdf"))
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
		assert.Contains(t, resp.Message, "PDF")
	})

	t.Run("File Too Large", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", "large.pdf")
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
