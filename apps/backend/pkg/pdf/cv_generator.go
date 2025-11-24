package pdf

import (
	"fmt"
	"strings"
	"time"

	"github.com/credfolio/apps/backend/internal/domain"
	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/pagesize"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

// CVGenerator handles PDF CV generation from domain models
type CVGenerator struct{}

// NewCVGenerator creates a new CV generator
func NewCVGenerator() *CVGenerator {
	return &CVGenerator{}
}

// GenerateCVFromProfile generates a PDF CV from a Profile domain model
// If jobMatch is provided, it will use the tailored summary and emphasize relevant content
func (g *CVGenerator) GenerateCVFromProfile(profile *domain.Profile, jobMatch *domain.JobMatch) ([]byte, error) {
	if profile == nil {
		return nil, fmt.Errorf("profile cannot be nil")
	}

	cfg := config.NewBuilder().
		WithPageSize(pagesize.A4).
		Build()

	mrt := maroto.New(cfg)

	// Use tailored summary if jobMatch is provided, otherwise use profile summary
	summary := profile.Summary
	if jobMatch != nil && jobMatch.TailoredSummary != "" {
		summary = jobMatch.TailoredSummary
	}

	// Header Section
	if summary != "" {
		mrt.AddRows(
			row.New(8).Add(
				col.New(12).Add(
					text.New("Professional Summary", props.Text{Size: 16}),
				),
			),
			row.New(10).Add(
				col.New(12).Add(
					text.New(summary, props.Text{Size: 10}),
				),
			),
		)
	}

	// Work Experience Section
	if len(profile.WorkExperiences) > 0 {
		mrt.AddRows(
			row.New(8).Add(
				col.New(12).Add(
					text.New("Work Experience", props.Text{Size: 16}),
				),
			),
		)

		for _, exp := range profile.WorkExperiences {
			g.addWorkExperience(mrt, exp, jobMatch)
		}
	}

	// Skills Section
	if len(profile.Skills) > 0 {
		skillNames := make([]string, len(profile.Skills))
		for i, skill := range profile.Skills {
			skillNames[i] = skill.Name
		}
		skillsText := strings.Join(skillNames, ", ")

		mrt.AddRows(
			row.New(8).Add(
				col.New(12).Add(
					text.New("Skills", props.Text{Size: 16}),
				),
			),
			row.New(6).Add(
				col.New(12).Add(
					text.New(skillsText, props.Text{Size: 10}),
				),
			),
		)
	}

	document, err := mrt.Generate()
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return document.GetBytes(), nil
}

// addWorkExperience adds a work experience entry to the PDF
func (g *CVGenerator) addWorkExperience(mrt core.Maroto, exp domain.WorkExperience, jobMatch *domain.JobMatch) {
	// Format dates
	startDate := formatDatePDF(exp.StartDate)
	endDate := "Present"
	if exp.EndDate != nil {
		endDate = formatDatePDF(*exp.EndDate)
	}
	dateRange := fmt.Sprintf("%s - %s", startDate, endDate)

	// Company and role
	companyRoleText := fmt.Sprintf("%s | %s", exp.Role, exp.CompanyName)

	mrt.AddRows(
		row.New(6).Add(
			col.New(8).Add(
				text.New(companyRoleText, props.Text{Size: 12}),
			),
			col.New(4).Add(
				text.New(dateRange, props.Text{Size: 9}),
			),
		),
	)

	// Description
	if exp.Description != "" {
		mrt.AddRows(
			row.New(5).Add(
				col.New(12).Add(
					text.New(exp.Description, props.Text{Size: 10}),
				),
			),
		)
	}

	// Credibility Highlights
	if len(exp.CredibilityHighlights) > 0 {
		mrt.AddRows(
			row.New(4).Add(
				col.New(12).Add(
					text.New("Employer Feedback:", props.Text{Size: 9}),
				),
			),
		)

		for _, highlight := range exp.CredibilityHighlights {
			quote := fmt.Sprintf("\"%s\"", highlight.Quote)
			mrt.AddRows(
				row.New(4).Add(
					col.New(12).Add(
						text.New(quote, props.Text{Size: 9}),
					),
				),
			)
		}
	}

	// Add spacing between work experiences
	mrt.AddRows(row.New(3).Add(col.New(12).Add(text.New("", props.Text{}))))
}

// formatDatePDF formats a date for PDF display
func formatDatePDF(t time.Time) string {
	return t.Format("Jan 2006")
}
