package generator

import (
	"bytes"
	"fmt"

	"github.com/credfolio/apps/backend/src/models"
	"github.com/go-pdf/fpdf"
)

type TailoringConfig struct {
	Summary           string
	HighlightSkillIDs []string
}

type CVGenerator struct{}

func NewCVGenerator() *CVGenerator {
	return &CVGenerator{}
}

func (g *CVGenerator) Generate(profile *models.UserProfile, tailoring *TailoringConfig) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 24)

	// Header
	pdf.Cell(0, 10, profile.FullName)
	pdf.Ln(12)
	pdf.SetFont("Arial", "", 12)
	pdf.SetTextColor(100, 100, 100)
	pdf.Cell(0, 10, profile.Email)
	pdf.Ln(12)
	pdf.SetTextColor(0, 0, 0)

	// Tailored Summary
	if tailoring != nil && tailoring.Summary != "" {
		pdf.SetFont("Arial", "I", 11)
		pdf.MultiCell(0, 5, tailoring.Summary, "", "", false)
		pdf.Ln(10)
	} else {
		pdf.Ln(8) // Standard spacing
	}

	// Experience Section
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Professional Experience")
	pdf.Line(10, pdf.GetY()+8, 200, pdf.GetY()+8)
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 12)
	for _, company := range profile.Companies {
		// Company Header
		pdf.SetFont("Arial", "B", 14)
		pdf.Cell(0, 8, company.Name)
		pdf.Ln(8)

		// Dates
		pdf.SetFont("Arial", "I", 10)
		startDate := company.StartDate.Format("Jan 2006")
		endDate := "Present"
		if company.EndDate != nil {
			endDate = company.EndDate.Format("Jan 2006")
		}
		pdf.Cell(0, 5, fmt.Sprintf("%s - %s", startDate, endDate))
		pdf.Ln(8)

		for _, role := range company.Roles {
			pdf.SetFont("Arial", "B", 12)
			pdf.Cell(0, 6, role.Title)
			pdf.Ln(7)

			pdf.SetFont("Arial", "", 11)
			pdf.MultiCell(0, 5, role.Description, "", "", false)
			pdf.Ln(6)

			if role.IsVerified {
				pdf.SetTextColor(0, 100, 0)
				pdf.SetFont("Arial", "I", 9)
				pdf.Cell(0, 5, "Verified by Reference")
				pdf.SetTextColor(0, 0, 0)
				pdf.Ln(6)
			}
		}
		pdf.Ln(5)
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
