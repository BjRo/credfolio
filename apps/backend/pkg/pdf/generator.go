package pdf

import (
	"fmt"

	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/pagesize"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

// Generator handles PDF generation for CVs
type Generator struct{}

// NewGenerator creates a new PDF generator
func NewGenerator() *Generator {
	return &Generator{}
}

// GenerateCV generates a PDF CV from profile data
func (g *Generator) GenerateCV(profileData *CVData) ([]byte, error) {
	cfg := config.NewBuilder().
		WithPageSize(pagesize.A4).
		Build()

	mrt := maroto.New(cfg)

	// Add title
	mrt.AddRows(
		row.New(10).Add(
			col.New(12).Add(
				text.New("CV", props.Text{Size: 20}),
			),
		),
	)

	if profileData.Summary != "" {
		mrt.AddRows(
			row.New(5).Add(
				col.New(12).Add(
					text.New("Summary", props.Text{Size: 16}),
				),
			),
			row.New(8).Add(
				col.New(12).Add(
					text.New(profileData.Summary, props.Text{Size: 10}),
				),
			),
		)
	}

	if len(profileData.WorkExperiences) > 0 {
		mrt.AddRows(
			row.New(10).Add(
				col.New(12).Add(
					text.New("Work Experience", props.Text{Size: 16}),
				),
			),
		)
		for _, exp := range profileData.WorkExperiences {
			title := fmt.Sprintf("%s - %s", exp.CompanyName, exp.Role)
			mrt.AddRows(
				row.New(8).Add(
					col.New(12).Add(
						text.New(title, props.Text{Size: 12}),
					),
				),
				row.New(6).Add(
					col.New(12).Add(
						text.New(exp.Description, props.Text{Size: 10}),
					),
				),
			)
		}
	}

	if len(profileData.Skills) > 0 {
		skillsText := ""
		for i, skill := range profileData.Skills {
			if i > 0 {
				skillsText += ", "
			}
			skillsText += skill
		}
		mrt.AddRows(
			row.New(10).Add(
				col.New(12).Add(
					text.New("Skills", props.Text{Size: 16}),
				),
			),
			row.New(8).Add(
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

// CVData represents the data structure for CV generation
type CVData struct {
	Summary         string
	WorkExperiences []WorkExperienceData
	Skills          []string
}

// WorkExperienceData represents work experience data for CV
type WorkExperienceData struct {
	CompanyName string
	Role        string
	Description string
}
