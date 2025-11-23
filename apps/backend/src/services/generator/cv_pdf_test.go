package generator

import (
	"testing"
	"time"

	"github.com/credfolio/apps/backend/src/models"
	"github.com/google/uuid"
)

func TestGenerate(t *testing.T) {
	generator := NewCVGenerator()

	profile := &models.UserProfile{
		ID:       uuid.New(),
		FullName: "John Doe",
		Email:    "john@example.com",
		Companies: []models.CompanyEntry{
			{
				Name:      "Acme",
				StartDate: time.Now(),
				Roles: []models.WorkExperience{
					{
						Title:       "Dev",
						Description: "Coding",
						IsVerified:  true,
					},
				},
			},
		},
	}

	pdfBytes, err := generator.Generate(profile, nil)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if len(pdfBytes) == 0 {
		t.Error("generated pdf is empty")
	}

	// Check header
	if string(pdfBytes[:4]) != "%PDF" {
		t.Errorf("invalid pdf header: %s", string(pdfBytes[:4]))
	}
}

func TestGenerateTailored(t *testing.T) {
	generator := NewCVGenerator()

	profile := &models.UserProfile{
		ID:       uuid.New(),
		FullName: "John Doe",
		Email:    "john@example.com",
		Companies: []models.CompanyEntry{
			{
				Name:      "Acme",
				StartDate: time.Now(),
				Roles: []models.WorkExperience{
					{
						Title:       "Dev",
						Description: "Coding",
						IsVerified:  true,
					},
				},
			},
		},
	}

	tailoring := &TailoringConfig{
		Summary: "This is a tailored summary.",
	}

	pdfBytes, err := generator.Generate(profile, tailoring)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if len(pdfBytes) == 0 {
		t.Error("generated pdf is empty")
	}
}
