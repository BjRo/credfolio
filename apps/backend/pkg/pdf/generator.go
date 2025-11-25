package pdf

import (
	"bytes"
	"errors"

	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

// ProfileData represents the data needed to generate a CV PDF
type ProfileData struct {
	Name            string
	Summary         string
	WorkExperiences []WorkExperienceData
	Skills          []string
	Highlights      []HighlightData
}

// WorkExperienceData represents work experience for PDF generation
type WorkExperienceData struct {
	CompanyName string
	Role        string
	StartDate   string
	EndDate     string
	Description string
	Highlights  []HighlightData
}

// HighlightData represents a credibility highlight for PDF generation
type HighlightData struct {
	Quote     string
	Sentiment string
}

// Generator handles PDF generation
type Generator struct{}

// NewGenerator creates a new PDF generator
func NewGenerator() *Generator {
	return &Generator{}
}

// GenerateCV generates a CV PDF from profile data
func (g *Generator) GenerateCV(data *ProfileData) ([]byte, error) {
	if data == nil {
		return nil, errors.New("profile data is required")
	}
	if data.Name == "" {
		return nil, errors.New("name is required")
	}

	m := maroto.New()

	// Add header with name
	m.AddRow(15,
		text.NewCol(12, data.Name, props.Text{
			Size:  18,
			Style: fontstyle.Bold,
			Align: align.Center,
		}),
	)

	// Add summary section
	if data.Summary != "" {
		g.addSection(m, "Summary")
		m.AddRow(10,
			text.NewCol(12, data.Summary, props.Text{
				Size: 10,
			}),
		)
	}

	// Add work experience section
	if len(data.WorkExperiences) > 0 {
		g.addSection(m, "Work Experience")
		for _, exp := range data.WorkExperiences {
			g.addWorkExperience(m, exp)
		}
	}

	// Add skills section
	if len(data.Skills) > 0 {
		g.addSection(m, "Skills")
		skillsText := ""
		for i, skill := range data.Skills {
			if i > 0 {
				skillsText += ", "
			}
			skillsText += skill
		}
		m.AddRow(8,
			text.NewCol(12, skillsText, props.Text{
				Size: 10,
			}),
		)
	}

	// Add credibility highlights section
	if len(data.Highlights) > 0 {
		g.addSection(m, "Credibility Highlights")
		for _, highlight := range data.Highlights {
			m.AddRow(8,
				text.NewCol(12, "\""+highlight.Quote+"\"", props.Text{
					Size:  10,
					Style: fontstyle.Italic,
				}),
			)
		}
	}

	// Generate PDF
	doc, err := m.Generate()
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	buf.Write(doc.GetBytes())

	return buf.Bytes(), nil
}

// addSection adds a section header to the PDF
func (g *Generator) addSection(m core.Maroto, title string) {
	m.AddRow(5)
	m.AddRow(10,
		text.NewCol(12, title, props.Text{
			Size:  14,
			Style: fontstyle.Bold,
			Top:   3,
		}),
	)
}

// addWorkExperience adds a work experience entry to the PDF
func (g *Generator) addWorkExperience(m core.Maroto, exp WorkExperienceData) {
	// Company and role
	dateRange := exp.StartDate
	if exp.EndDate != "" {
		dateRange += " - " + exp.EndDate
	} else {
		dateRange += " - Present"
	}

	m.AddRow(8,
		text.NewCol(8, exp.CompanyName+" - "+exp.Role, props.Text{
			Size:  11,
			Style: fontstyle.Bold,
		}),
		text.NewCol(4, dateRange, props.Text{
			Size:  10,
			Align: align.Right,
		}),
	)

	// Description
	if exp.Description != "" {
		m.AddRow(6,
			text.NewCol(12, exp.Description, props.Text{
				Size: 10,
			}),
		)
	}

	// Experience-specific highlights
	for _, highlight := range exp.Highlights {
		m.AddRow(6,
			text.NewCol(12, "• \""+highlight.Quote+"\"", props.Text{
				Size:  9,
				Style: fontstyle.Italic,
			}),
		)
	}

	m.AddRow(3)
}
