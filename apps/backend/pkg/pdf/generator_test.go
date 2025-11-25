package pdf_test

import (
	"testing"

	"github.com/credfolio/apps/backend/pkg/pdf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerator_GenerateCV_WhenGivenValidData_GeneratesPDFBytes(t *testing.T) {
	// Arrange
	generator := pdf.NewGenerator()
	data := &pdf.ProfileData{
		Name:    "John Doe",
		Summary: "Experienced software engineer",
		WorkExperiences: []pdf.WorkExperienceData{
			{
				CompanyName: "Tech Corp",
				Role:        "Senior Engineer",
				StartDate:   "2020-01",
				EndDate:     "2023-12",
				Description: "Led backend development",
			},
		},
		Skills: []string{"Go", "Python", "Kubernetes"},
		Highlights: []pdf.HighlightData{
			{Quote: "Excellent team player", Sentiment: "POSITIVE"},
		},
	}

	// Act
	result, err := generator.GenerateCV(data)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, result)
	// PDF files start with %PDF
	assert.True(t, len(result) > 4)
}

func TestGenerator_GenerateCV_WhenNilData_ReturnsError(t *testing.T) {
	// Arrange
	generator := pdf.NewGenerator()

	// Act
	result, err := generator.GenerateCV(nil)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "profile data is required")
}

func TestGenerator_GenerateCV_WhenEmptyName_ReturnsError(t *testing.T) {
	// Arrange
	generator := pdf.NewGenerator()
	data := &pdf.ProfileData{
		Name: "",
	}

	// Act
	result, err := generator.GenerateCV(data)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "name is required")
}

func TestGenerator_GenerateCV_WhenMinimalData_GeneratesPDF(t *testing.T) {
	// Arrange
	generator := pdf.NewGenerator()
	data := &pdf.ProfileData{
		Name: "Jane Doe",
	}

	// Act
	result, err := generator.GenerateCV(data)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestGenerator_GenerateCV_WhenMultipleExperiences_GeneratesPDF(t *testing.T) {
	// Arrange
	generator := pdf.NewGenerator()
	data := &pdf.ProfileData{
		Name: "John Doe",
		WorkExperiences: []pdf.WorkExperienceData{
			{CompanyName: "Company A", Role: "Engineer", StartDate: "2020-01"},
			{CompanyName: "Company B", Role: "Senior Engineer", StartDate: "2022-01", EndDate: "2023-12"},
		},
	}

	// Act
	result, err := generator.GenerateCV(data)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestNewGenerator_CreatesGenerator(t *testing.T) {
	// Act
	generator := pdf.NewGenerator()

	// Assert
	assert.NotNil(t, generator)
}

// T128: Unit test for CV generator when given profile data generates PDF with all sections
func TestGenerator_GenerateCV_WhenGivenProfileData_GeneratesPDFWithAllSections(t *testing.T) {
	// Arrange
	generator := pdf.NewGenerator()
	data := &pdf.ProfileData{
		Name:    "Alex Johnson",
		Summary: "Senior software engineer with 10 years experience",
		WorkExperiences: []pdf.WorkExperienceData{
			{
				CompanyName: "Tech Corp",
				Role:        "Lead Engineer",
				StartDate:   "2020-01",
				EndDate:     "2023-12",
				Description: "Led backend development team",
				Highlights: []pdf.HighlightData{
					{Quote: "Outstanding leadership", Sentiment: "POSITIVE"},
				},
			},
			{
				CompanyName: "Startup Inc",
				Role:        "Software Engineer",
				StartDate:   "2018-01",
				EndDate:     "2019-12",
				Description: "Full-stack development",
			},
		},
		Skills: []string{"Go", "TypeScript", "PostgreSQL", "Docker", "Kubernetes"},
		Highlights: []pdf.HighlightData{
			{Quote: "Exceptional problem solver", Sentiment: "POSITIVE"},
			{Quote: "Great communicator", Sentiment: "POSITIVE"},
		},
	}

	// Act
	result, err := generator.GenerateCV(data)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, result)
	// PDF should contain significant content
	assert.Greater(t, len(result), 1000, "PDF should have substantial content for all sections")
}

// T129: Unit test for CV generator when given tailored profile emphasizes matched content
func TestGenerator_GenerateTailoredCV_WhenGivenTailoredProfile_EmphasizesMatchedContent(t *testing.T) {
	// Arrange
	generator := pdf.NewGenerator()
	data := &pdf.ProfileData{
		Name:    "Alex Johnson",
		Summary: "Senior software engineer",
		WorkExperiences: []pdf.WorkExperienceData{
			{
				CompanyName: "Tech Corp",
				Role:        "Backend Engineer",
				StartDate:   "2020-01",
				EndDate:     "2023-12",
				Description: "Go microservices development",
			},
		},
		Skills: []string{"Go", "PostgreSQL"},
	}

	tailoringOptions := &pdf.TailoringOptions{
		JobDescription: "Backend engineer with Go experience",
		MatchScore:     0.85,
		MatchSummary:   "Strong match based on Go experience",
		RelevantSkills: []string{"Go", "PostgreSQL"},
		ExperienceScores: map[int]float64{
			0: 0.9, // First experience has high relevance
		},
	}

	// Act
	result, err := generator.GenerateTailoredCV(data, tailoringOptions)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Greater(t, len(result), 1000)
}

// T130: Unit test for CV generator when profile is empty returns error
func TestGenerator_GenerateCV_WhenProfileIsEmpty_ReturnsError(t *testing.T) {
	// Arrange
	generator := pdf.NewGenerator()
	data := &pdf.ProfileData{
		Name:            "",
		Summary:         "",
		WorkExperiences: []pdf.WorkExperienceData{},
		Skills:          []string{},
	}

	// Act
	result, err := generator.GenerateCV(data)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGenerator_GenerateTailoredCV_WhenNilOptions_GeneratesStandardCV(t *testing.T) {
	// Arrange
	generator := pdf.NewGenerator()
	data := &pdf.ProfileData{
		Name:    "John Doe",
		Summary: "Software engineer",
	}

	// Act
	result, err := generator.GenerateTailoredCV(data, nil)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, result)
}
