package pdf

import (
	"bytes"
	"errors"
	"fmt"

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

// TailoringOptions contains options for tailored CV generation
type TailoringOptions struct {
	JobDescription   string
	MatchScore       float64
	MatchSummary     string
	RelevantSkills   []string
	ExperienceScores map[int]float64 // Index -> relevance score
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

// GenerateTailoredCV generates a CV PDF with tailored content emphasis
func (g *Generator) GenerateTailoredCV(data *ProfileData, options *TailoringOptions) ([]byte, error) {
	if data == nil {
		return nil, errors.New("profile data is required")
	}
	if data.Name == "" {
		return nil, errors.New("name is required")
	}

	// If no tailoring options, generate standard CV
	if options == nil {
		return g.GenerateCV(data)
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

	// Add match score banner
	matchText := ""
	switch {
	case options.MatchScore >= 0.7:
		matchText = "Strong Match (" + formatPercent(options.MatchScore) + ")"
	case options.MatchScore >= 0.4:
		matchText = "Moderate Match (" + formatPercent(options.MatchScore) + ")"
	default:
		matchText = "Match Score: " + formatPercent(options.MatchScore)
	}
	m.AddRow(10,
		text.NewCol(12, matchText, props.Text{
			Size:  12,
			Style: fontstyle.BoldItalic,
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

	// Add match summary if available
	if options.MatchSummary != "" {
		m.AddRow(8,
			text.NewCol(12, options.MatchSummary, props.Text{
				Size:  10,
				Style: fontstyle.Italic,
			}),
		)
	}

	// Add work experience section with relevance indicators
	if len(data.WorkExperiences) > 0 {
		g.addSection(m, "Work Experience")
		for i, exp := range data.WorkExperiences {
			score, hasScore := options.ExperienceScores[i]
			g.addTailoredWorkExperience(m, exp, score, hasScore)
		}
	}

	// Add skills section with relevant skills highlighted
	if len(data.Skills) > 0 {
		g.addSection(m, "Skills")

		// Show relevant skills first
		relevantSet := make(map[string]bool)
		for _, skill := range options.RelevantSkills {
			relevantSet[skill] = true
		}

		relevantSkills := ""
		otherSkills := ""

		for _, skill := range data.Skills {
			if relevantSet[skill] {
				if relevantSkills != "" {
					relevantSkills += ", "
				}
				relevantSkills += skill + " ★"
			} else {
				if otherSkills != "" {
					otherSkills += ", "
				}
				otherSkills += skill
			}
		}

		if relevantSkills != "" {
			m.AddRow(8,
				text.NewCol(12, "Relevant: "+relevantSkills, props.Text{
					Size:  10,
					Style: fontstyle.Bold,
				}),
			)
		}
		if otherSkills != "" {
			m.AddRow(8,
				text.NewCol(12, "Other: "+otherSkills, props.Text{
					Size: 10,
				}),
			)
		}
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

// formatPercent formats a float as a percentage string
func formatPercent(f float64) string {
	percent := int(f * 100)
	return fmt.Sprintf("%d%%", percent)
}

// addTailoredWorkExperience adds a work experience entry with relevance indicator
func (g *Generator) addTailoredWorkExperience(m core.Maroto, exp WorkExperienceData, score float64, hasScore bool) {
	// Company and role
	dateRange := exp.StartDate
	if exp.EndDate != "" {
		dateRange += " - " + exp.EndDate
	} else {
		dateRange += " - Present"
	}

	// Add relevance indicator if available
	relevanceText := ""
	if hasScore {
		if score >= 0.7 {
			relevanceText = " [High Relevance]"
		} else if score >= 0.4 {
			relevanceText = " [Moderate Relevance]"
		}
	}

	m.AddRow(8,
		text.NewCol(8, exp.CompanyName+" - "+exp.Role+relevanceText, props.Text{
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
