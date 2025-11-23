package profile

import (
	"context"
	"strings"
	"testing"

	"github.com/credfolio/apps/backend/src/models"
	"github.com/credfolio/apps/backend/src/services/storage"
	"github.com/google/uuid"
)

type MockRepo struct{}

func (m *MockRepo) CreateReferenceLetter(ctx context.Context, letter *models.ReferenceLetter) error {
	return nil
}
func (m *MockRepo) UpsertCompany(ctx context.Context, company *models.CompanyEntry) error { return nil }
func (m *MockRepo) CreateWorkExperience(ctx context.Context, experience *models.WorkExperience) error {
	return nil
}
func (m *MockRepo) UpsertSkill(ctx context.Context, skill *models.Skill) error { return nil }
func (m *MockRepo) LinkExperienceSkill(ctx context.Context, experienceID, skillID uuid.UUID) error {
	return nil
}

func (m *MockRepo) GetUser(ctx context.Context, id uuid.UUID) (*models.UserProfile, error) {
	return &models.UserProfile{ID: id, FullName: "Test User"}, nil
}

func TestProcessUpload(t *testing.T) {
	repo := &MockRepo{}

	// Mock Extractor dependencies (reusing types from extractor_test.go if in same package)
	// Since it's same package 'profile', we can reuse MockPDFExtractor and MockLLMClient if they are exported or in same package scope.
	// Wait, they are in extractor_test.go which is usually only compiled for tests of that file?
	// No, `go test` compiles all _test.go files in package. So they are available.

	mockPDF := &MockPDFExtractor{Text: "dummy text"}
	mockLLM := &MockLLMClient{
		Response: `{"full_name": "Test", "company": {"name":"TestCo", "start_date":"2020-01-01"}, "role": {"title":"Dev", "skills":[]}}`,
	}
	extractor := NewExtractor(mockPDF, mockLLM)

	// Real storage in temp
	tmpDir := t.TempDir()
	store, _ := storage.NewLocalStorage(tmpDir)

	svc := NewService(repo, extractor, store, mockLLM)

	err := svc.ProcessUpload(context.Background(), uuid.New(), "test.pdf", strings.NewReader("dummy content"))
	if err != nil {
		t.Fatalf("ProcessUpload failed: %v", err)
	}
}

