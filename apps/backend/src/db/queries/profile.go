package queries

import (
	"context"

	"github.com/credfolio/apps/backend/src/db"
	"github.com/credfolio/apps/backend/src/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type PgxPool interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

type ProfileQueries struct {
	Pool PgxPool
}

func NewProfileQueries(d *db.DB) *ProfileQueries {
	return &ProfileQueries{Pool: d.Pool}
}

func (q *ProfileQueries) CreateReferenceLetter(ctx context.Context, letter *models.ReferenceLetter) error {
	sql := `INSERT INTO reference_letters (id, user_id, filename, storage_path, content_hash, upload_date, extraction_status, extracted_metadata)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := q.Pool.Exec(ctx, sql,
		letter.ID, letter.UserID, letter.Filename, letter.StoragePath,
		letter.ContentHash, letter.UploadDate, letter.ExtractionStatus, letter.ExtractedMetadata)
	return err
}

func (q *ProfileQueries) UpsertCompany(ctx context.Context, company *models.CompanyEntry) error {
	sql := `SELECT id FROM companies WHERE user_id = $1 AND name = $2`
	var existingID uuid.UUID
	err := q.Pool.QueryRow(ctx, sql, company.UserID, company.Name).Scan(&existingID)

	if err == nil {
		company.ID = existingID
		return nil
	} else if err != pgx.ErrNoRows {
		return err
	}

	sql = `INSERT INTO companies (id, user_id, name, logo_url, start_date, end_date)
           VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = q.Pool.Exec(ctx, sql, company.ID, company.UserID, company.Name, company.LogoURL, company.StartDate, company.EndDate)
	return err
}

func (q *ProfileQueries) CreateWorkExperience(ctx context.Context, experience *models.WorkExperience) error {
	sql := `INSERT INTO work_experiences (id, company_id, title, start_date, end_date, description, source, employer_feedback, reference_letter_id, is_verified)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := q.Pool.Exec(ctx, sql,
		experience.ID, experience.CompanyID, experience.Title,
		experience.StartDate, experience.EndDate, experience.Description,
		experience.Source, experience.EmployerFeedback,
		experience.ReferenceLetterID, experience.IsVerified)
	return err
}

func (q *ProfileQueries) UpsertSkill(ctx context.Context, skill *models.Skill) error {
	sql := `SELECT id FROM skills WHERE name = $1`
	var existingID uuid.UUID
	err := q.Pool.QueryRow(ctx, sql, skill.Name).Scan(&existingID)

	if err == nil {
		skill.ID = existingID
		return nil
	} else if err != pgx.ErrNoRows {
		return err
	}

	sql = `INSERT INTO skills (id, name) VALUES ($1, $2)`
	_, err = q.Pool.Exec(ctx, sql, skill.ID, skill.Name)
	return err
}

func (q *ProfileQueries) LinkExperienceSkill(ctx context.Context, experienceID, skillID uuid.UUID) error {
	sql := `INSERT INTO experience_skills (work_experience_id, skill_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := q.Pool.Exec(ctx, sql, experienceID, skillID)
	return err
}

func (q *ProfileQueries) GetUser(ctx context.Context, id uuid.UUID) (*models.UserProfile, error) {
	user := &models.UserProfile{}
	err := q.Pool.QueryRow(ctx, "SELECT id, email, full_name, created_at FROM users WHERE id = $1", id).Scan(
		&user.ID, &user.Email, &user.FullName, &user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	companies, err := q.getCompaniesForUser(ctx, id)
	if err != nil {
		return nil, err
	}
	user.Companies = companies

	return user, nil
}

func (q *ProfileQueries) getCompaniesForUser(ctx context.Context, userID uuid.UUID) ([]models.CompanyEntry, error) {
	rows, err := q.Pool.Query(ctx, "SELECT id, name, logo_url, start_date, end_date FROM companies WHERE user_id = $1 ORDER BY start_date DESC", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companies []models.CompanyEntry
	for rows.Next() {
		var c models.CompanyEntry
		c.UserID = userID
		if err := rows.Scan(&c.ID, &c.Name, &c.LogoURL, &c.StartDate, &c.EndDate); err != nil {
			return nil, err
		}

		roles, err := q.getRolesForCompany(ctx, c.ID)
		if err != nil {
			return nil, err
		}
		c.Roles = roles
		companies = append(companies, c)
	}
	return companies, nil
}

func (q *ProfileQueries) getRolesForCompany(ctx context.Context, companyID uuid.UUID) ([]models.WorkExperience, error) {
	rows, err := q.Pool.Query(ctx, "SELECT id, title, start_date, end_date, description, source, employer_feedback, is_verified FROM work_experiences WHERE company_id = $1 ORDER BY start_date DESC", companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []models.WorkExperience
	for rows.Next() {
		var r models.WorkExperience
		r.CompanyID = companyID
		if err := rows.Scan(&r.ID, &r.Title, &r.StartDate, &r.EndDate, &r.Description, &r.Source, &r.EmployerFeedback, &r.IsVerified); err != nil {
			return nil, err
		}

		skills, err := q.getSkillsForExperience(ctx, r.ID)
		if err != nil {
			return nil, err
		}
		r.Skills = skills
		roles = append(roles, r)
	}
	return roles, nil
}

func (q *ProfileQueries) getSkillsForExperience(ctx context.Context, expID uuid.UUID) ([]models.Skill, error) {
	sql := `SELECT s.id, s.name FROM skills s JOIN experience_skills es ON s.id = es.skill_id WHERE es.work_experience_id = $1`
	rows, err := q.Pool.Query(ctx, sql, expID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var skills []models.Skill
	for rows.Next() {
		var s models.Skill
		if err := rows.Scan(&s.ID, &s.Name); err != nil {
			return nil, err
		}
		skills = append(skills, s)
	}
	return skills, nil
}
