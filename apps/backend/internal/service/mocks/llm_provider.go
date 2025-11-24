package mocks

import (
	"context"

	"github.com/credfolio/apps/backend/internal/service"
	"github.com/stretchr/testify/mock"
)

type MockLLMProvider struct {
	mock.Mock
}

func (m *MockLLMProvider) ExtractProfileData(ctx context.Context, text string) (*service.ExtractedProfileData, error) {
	args := m.Called(ctx, text)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.ExtractedProfileData), args.Error(1)
}

func (m *MockLLMProvider) ExtractCredibility(ctx context.Context, text string) (*service.CredibilityData, error) {
	args := m.Called(ctx, text)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.CredibilityData), args.Error(1)
}

func (m *MockLLMProvider) TailorProfile(ctx context.Context, profileText string, jobDescription string) (string, float64, error) {
	args := m.Called(ctx, profileText, jobDescription)
	return args.String(0), args.Get(1).(float64), args.Error(2)
}
