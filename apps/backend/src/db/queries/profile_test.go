package queries

import (
	"context"
	"testing"
	"time"

	"github.com/credfolio/apps/backend/src/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
)

func TestUpsertCompany(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	q := &ProfileQueries{
		Pool: mock,
	}

	userID := uuid.New()
	company := &models.CompanyEntry{
		ID:        uuid.New(),
		UserID:    userID,
		Name:      "Test Corp",
		StartDate: time.Now(),
	}

	// Expect check for existing company
	mock.ExpectQuery("SELECT id FROM companies").
		WithArgs(userID, "Test Corp").
		WillReturnError(pgx.ErrNoRows)

	// Expect insert
	mock.ExpectExec("INSERT INTO companies").
		WithArgs(company.ID, userID, "Test Corp", company.LogoURL, company.StartDate, company.EndDate).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err = q.UpsertCompany(context.Background(), company)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
