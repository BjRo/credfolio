package mocks

import (
	"context"

	"github.com/credfolio/apps/backend/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockProfileRepository struct {
	mock.Mock
}

func (m *MockProfileRepository) Create(ctx context.Context, profile *domain.Profile) error {
	args := m.Called(ctx, profile)
	return args.Error(0)
}

func (m *MockProfileRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Profile, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Profile), args.Error(1)
}

func (m *MockProfileRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Profile, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Profile), args.Error(1)
}

func (m *MockProfileRepository) Update(ctx context.Context, profile *domain.Profile) error {
	args := m.Called(ctx, profile)
	return args.Error(0)
}

func (m *MockProfileRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockWorkExperienceRepository struct {
	mock.Mock
}

func (m *MockWorkExperienceRepository) Create(ctx context.Context, we *domain.WorkExperience) error {
	args := m.Called(ctx, we)
	return args.Error(0)
}

func (m *MockWorkExperienceRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.WorkExperience, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.WorkExperience), args.Error(1)
}

func (m *MockWorkExperienceRepository) GetByProfileID(ctx context.Context, profileID uuid.UUID) ([]*domain.WorkExperience, error) {
	args := m.Called(ctx, profileID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.WorkExperience), args.Error(1)
}

func (m *MockWorkExperienceRepository) Update(ctx context.Context, we *domain.WorkExperience) error {
	args := m.Called(ctx, we)
	return args.Error(0)
}

func (m *MockWorkExperienceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockCredibilityHighlightRepository struct {
	mock.Mock
}

func (m *MockCredibilityHighlightRepository) Create(ctx context.Context, ch *domain.CredibilityHighlight) error {
	args := m.Called(ctx, ch)
	return args.Error(0)
}

func (m *MockCredibilityHighlightRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.CredibilityHighlight, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.CredibilityHighlight), args.Error(1)
}

func (m *MockCredibilityHighlightRepository) GetByWorkExperienceID(ctx context.Context, weID uuid.UUID) ([]*domain.CredibilityHighlight, error) {
	args := m.Called(ctx, weID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.CredibilityHighlight), args.Error(1)
}

func (m *MockCredibilityHighlightRepository) Update(ctx context.Context, ch *domain.CredibilityHighlight) error {
	args := m.Called(ctx, ch)
	return args.Error(0)
}

func (m *MockCredibilityHighlightRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockReferenceLetterRepository struct {
	mock.Mock
}

func (m *MockReferenceLetterRepository) Create(ctx context.Context, letter *domain.ReferenceLetter) error {
	args := m.Called(ctx, letter)
	return args.Error(0)
}

func (m *MockReferenceLetterRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.ReferenceLetter, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ReferenceLetter), args.Error(1)
}

func (m *MockReferenceLetterRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.ReferenceLetter, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.ReferenceLetter), args.Error(1)
}

func (m *MockReferenceLetterRepository) Update(ctx context.Context, letter *domain.ReferenceLetter) error {
	args := m.Called(ctx, letter)
	return args.Error(0)
}

func (m *MockReferenceLetterRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockJobMatchRepository struct {
	mock.Mock
}

func (m *MockJobMatchRepository) Create(ctx context.Context, jobMatch *domain.JobMatch) error {
	args := m.Called(ctx, jobMatch)
	return args.Error(0)
}

func (m *MockJobMatchRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.JobMatch, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.JobMatch), args.Error(1)
}

func (m *MockJobMatchRepository) GetByProfileID(ctx context.Context, profileID uuid.UUID) ([]*domain.JobMatch, error) {
	args := m.Called(ctx, profileID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.JobMatch), args.Error(1)
}

func (m *MockJobMatchRepository) Update(ctx context.Context, jobMatch *domain.JobMatch) error {
	args := m.Called(ctx, jobMatch)
	return args.Error(0)
}

func (m *MockJobMatchRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
