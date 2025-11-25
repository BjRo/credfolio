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
