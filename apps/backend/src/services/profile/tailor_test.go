package profile

import (
	"context"
	"testing"

	"github.com/credfolio/apps/backend/src/services/storage"
	"github.com/google/uuid"
)

func TestTailorProfile(t *testing.T) {
	repo := &MockRepo{} // MockRepo from service_test.go needs GetUser

	// Mock LLM for tailoring
	mockLLM := &MockLLMClient{
		Response: `{"relevant_skill_ids": ["uuid1"], "match_score": 90}`,
	}

	mockPDF := &MockPDFExtractor{}
	extractor := NewExtractor(mockPDF, mockLLM)

	tmpDir := t.TempDir()
	store, _ := storage.NewLocalStorage(tmpDir)

	svc := NewService(repo, extractor, store, mockLLM)

	res, err := svc.TailorProfile(context.Background(), uuid.New(), "Job Desc")
	if err != nil {
		t.Fatalf("TailorProfile failed: %v", err)
	}

	if res.MatchScore != 90 {
		t.Errorf("expected score 90, got %d", res.MatchScore)
	}
}
