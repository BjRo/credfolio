package profile

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/credfolio/apps/backend/src/models"
	"github.com/credfolio/apps/backend/src/services/storage"
	"github.com/google/uuid"
)

type Repository interface {
	CreateReferenceLetter(ctx context.Context, letter *models.ReferenceLetter) error
	// UpsertCompany checks if company exists for user by name. If so, updates ID in struct. If not, creates.
	UpsertCompany(ctx context.Context, company *models.CompanyEntry) error
	CreateWorkExperience(ctx context.Context, experience *models.WorkExperience) error
	UpsertSkill(ctx context.Context, skill *models.Skill) error
	LinkExperienceSkill(ctx context.Context, experienceID, skillID uuid.UUID) error
	GetUser(ctx context.Context, id uuid.UUID) (*models.UserProfile, error)
}

type Service struct {
	repo      Repository
	extractor *Extractor
	storage   *storage.LocalStorage
	llmClient LLMClient
}

func NewService(repo Repository, extractor *Extractor, storage *storage.LocalStorage, llmClient LLMClient) *Service {
	return &Service{
		repo:      repo,
		extractor: extractor,
		storage:   storage,
		llmClient: llmClient,
	}
}

func (s *Service) GetProfile(ctx context.Context, userID uuid.UUID) (*models.UserProfile, error) {
	return s.repo.GetUser(ctx, userID)
}

func (s *Service) ProcessUpload(ctx context.Context, userID uuid.UUID, filename string, content io.Reader) error {
	// 1. Save file
	path, err := s.storage.Save(filename, content)
	if err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	// 2. Open file for extraction
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open saved file: %w", err)
	}
	defer file.Close()

	stat, _ := file.Stat()
	fmt.Printf("DEBUG: Processing file %s, size: %d bytes\n", path, stat.Size())

	// 3. Extract
	data, err := s.extractor.Extract(ctx, file, stat.Size())
	if err != nil {
		return fmt.Errorf("extraction failed: %w", err)
	}

	// 4. Map to models and Save

	// Reference Letter
	refLetter := &models.ReferenceLetter{
		ID:               uuid.New(),
		UserID:           userID,
		Filename:         filename,
		StoragePath:      path,
		UploadDate:       time.Now(),
		ExtractionStatus: models.StatusCompleted,
	}
	if err := s.repo.CreateReferenceLetter(ctx, refLetter); err != nil {
		return fmt.Errorf("failed to save ref letter: %w", err)
	}

	// Company
	company := &models.CompanyEntry{
		ID:        uuid.New(),
		UserID:    userID,
		Name:      data.Company.Name,
		StartDate: parseDate(data.Company.StartDate),
		EndDate:   parseDatePtr(data.Company.EndDate),
	}
	if err := s.repo.UpsertCompany(ctx, company); err != nil {
		return fmt.Errorf("failed to upsert company: %w", err)
	}

	// Work Experience
	experience := &models.WorkExperience{
		ID:                uuid.New(),
		CompanyID:         company.ID,
		Title:             data.Role.Title,
		StartDate:         company.StartDate,
		EndDate:           company.EndDate,
		Description:       data.Role.Description,
		Source:            models.SourceVerified,
		EmployerFeedback:  data.Role.EmployerFeedback,
		ReferenceLetterID: &refLetter.ID,
		IsVerified:        true,
	}
	if err := s.repo.CreateWorkExperience(ctx, experience); err != nil {
		return fmt.Errorf("failed to create experience: %w", err)
	}

	// Skills
	for _, skillName := range data.Role.Skills {
		skill := &models.Skill{
			ID:   uuid.New(),
			Name: skillName,
		}
		if err := s.repo.UpsertSkill(ctx, skill); err != nil {
			continue
		}
		_ = s.repo.LinkExperienceSkill(ctx, experience.ID, skill.ID)
	}

	return nil
}

func parseDate(d string) time.Time {
	t, _ := time.Parse("2006-01-02", d)
	return t
}

func parseDatePtr(d *string) *time.Time {
	if d == nil {
		return nil
	}
	t := parseDate(*d)
	return &t
}

