package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/credfolio/apps/backend/internal/domain"
	"github.com/credfolio/apps/backend/internal/repository/mocks"
	"github.com/credfolio/apps/backend/internal/service"
	servicemocks "github.com/credfolio/apps/backend/internal/service/mocks"
	"github.com/credfolio/apps/backend/pkg/logger"
	pdfmocks "github.com/credfolio/apps/backend/pkg/pdf/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGenerateProfileFromReferences(t *testing.T) {
	mockProfileRepo := new(mocks.MockProfileRepository)
	mockWorkExpRepo := new(mocks.MockWorkExperienceRepository)
	mockCredibilityRepo := new(mocks.MockCredibilityHighlightRepository)
	mockRefLetterRepo := new(mocks.MockReferenceLetterRepository)
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

	userID := uuid.New()
	letterID := uuid.New()
	profileID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		// Setup Mocks
		mockProfileRepo.On("GetByUserID", mock.Anything, userID).Return(&domain.Profile{
			ID:     profileID,
			UserID: userID,
		}, nil)

		mockRefLetterRepo.On("GetByID", mock.Anything, letterID).Return(&domain.ReferenceLetter{
			ID:            letterID,
			UserID:        userID,
			ExtractedText: "Sample Reference Letter",
		}, nil)

		mockLLMProvider.On("ExtractProfileData", mock.Anything, "Sample Reference Letter").Return(&service.ExtractedProfileData{
			CompanyName: "Tech Corp",
			Role:        "Senior Developer",
			StartDate:   "2023-01-01",
			Description: "Great work",
		}, nil)

		mockLLMProvider.On("ExtractCredibility", mock.Anything, "Sample Reference Letter").Return(&service.CredibilityData{
			Quotes:    []string{"Best dev ever"},
			Sentiment: "POSITIVE",
		}, nil)

		mockWorkExpRepo.On("Create", mock.Anything, mock.MatchedBy(func(we *domain.WorkExperience) bool {
			return we.CompanyName == "Tech Corp" && we.Role == "Senior Developer" && we.ProfileID == profileID
		})).Return(nil)

		mockCredibilityRepo.On("Create", mock.Anything, mock.MatchedBy(func(ch *domain.CredibilityHighlight) bool {
			return ch.Quote == "Best dev ever" && ch.Sentiment == domain.SentimentPositive
		})).Return(nil)

		mockProfileRepo.On("GetByID", mock.Anything, profileID).Return(&domain.Profile{
			ID:     profileID,
			UserID: userID,
		}, nil)

		// Execute
		result, err := svc.GenerateProfileFromReferences(context.Background(), userID, []uuid.UUID{letterID})

		// Verify
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, profileID, result.ID)

		mockProfileRepo.AssertExpectations(t)
		mockRefLetterRepo.AssertExpectations(t)
		mockLLMProvider.AssertExpectations(t)
		mockWorkExpRepo.AssertExpectations(t)
		mockCredibilityRepo.AssertExpectations(t)
	})

	t.Run("Profile Creation on Not Found", func(t *testing.T) {
		// Reset Mocks
		mockProfileRepo = new(mocks.MockProfileRepository)
		mockRefLetterRepo = new(mocks.MockReferenceLetterRepository)
		mockLLMProvider = new(servicemocks.MockLLMProvider)
		mockWorkExpRepo = new(mocks.MockWorkExperienceRepository)
		mockCredibilityRepo = new(mocks.MockCredibilityHighlightRepository)
		mockPDFExtractor = new(pdfmocks.MockPDFExtractor)

		svc = service.NewProfileService(
			mockProfileRepo,
			mockWorkExpRepo,
			mockCredibilityRepo,
			mockRefLetterRepo,
			mockLLMProvider,
			mockPDFExtractor,
			appLogger,
		)

		// Setup Mocks
		mockProfileRepo.On("GetByUserID", mock.Anything, userID).Return(nil, errors.New("not found")).Once()

		// Expect Create Profile
		mockProfileRepo.On("Create", mock.Anything, mock.MatchedBy(func(p *domain.Profile) bool {
			return p.UserID == userID
		})).Return(nil).Run(func(args mock.Arguments) {
			p := args.Get(1).(*domain.Profile)
			p.ID = profileID // simulate DB setting ID
		})

		// Since Create is called, we proceed with normal flow...
		mockProfileRepo.On("GetByID", mock.Anything, profileID).Return(&domain.Profile{
			ID:     profileID,
			UserID: userID,
		}, nil)

		// Execute
		result, err := svc.GenerateProfileFromReferences(context.Background(), userID, []uuid.UUID{})

		// Verify
		assert.NoError(t, err)
		assert.NotNil(t, result)

		mockProfileRepo.AssertExpectations(t)
	})
}
