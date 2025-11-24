package ai

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/credfolio/apps/backend/internal/service"
	servicemocks "github.com/credfolio/apps/backend/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCachedLLMProvider_ExtractProfileData_WhenCacheMiss_ThenCallsProviderAndCaches(t *testing.T) {
	// Arrange
	mockProvider := new(servicemocks.MockLLMProvider)
	cachedProvider := NewCachedLLMProvider(mockProvider)
	ctx := context.Background()
	text := "Sample reference letter text"
	expectedData := &service.ExtractedProfileData{
		CompanyName:  "Test Company",
		Role:         "Software Engineer",
		StartDate:    "2020-01-01",
		EndDate:      "2022-12-31",
		Skills:       []string{"Go", "Python"},
		Achievements: []string{"Built API"},
		Description:  "Worked on backend systems",
	}

	mockProvider.On("ExtractProfileData", ctx, text).Return(expectedData, nil).Once()

	// Act - First call (cache miss)
	result1, err1 := cachedProvider.ExtractProfileData(ctx, text)

	// Assert - First call
	require.NoError(t, err1)
	assert.Equal(t, expectedData, result1)
	mockProvider.AssertExpectations(t)

	// Act - Second call (cache hit)
	result2, err2 := cachedProvider.ExtractProfileData(ctx, text)

	// Assert - Second call should return cached result without calling provider
	require.NoError(t, err2)
	assert.Equal(t, expectedData, result2)
	assert.Equal(t, result1, result2)  // Same pointer means it's cached
	mockProvider.AssertExpectations(t) // Should not have been called again
}

func TestCachedLLMProvider_ExtractProfileData_WhenProviderError_ThenReturnsErrorAndDoesNotCache(t *testing.T) {
	// Arrange
	mockProvider := new(servicemocks.MockLLMProvider)
	cachedProvider := NewCachedLLMProvider(mockProvider)
	ctx := context.Background()
	text := "Sample reference letter text"
	expectedError := errors.New("provider error")

	mockProvider.On("ExtractProfileData", ctx, text).Return(nil, expectedError).Twice()

	// Act - First call
	result1, err1 := cachedProvider.ExtractProfileData(ctx, text)

	// Assert - First call
	assert.Error(t, err1)
	assert.Nil(t, result1)
	assert.Equal(t, expectedError, err1)

	// Act - Second call (should still call provider since error wasn't cached)
	result2, err2 := cachedProvider.ExtractProfileData(ctx, text)

	// Assert - Second call
	assert.Error(t, err2)
	assert.Nil(t, result2)
	mockProvider.AssertExpectations(t)
}

func TestCachedLLMProvider_ExtractCredibility_WhenCacheMiss_ThenCallsProviderAndCaches(t *testing.T) {
	// Arrange
	mockProvider := new(servicemocks.MockLLMProvider)
	cachedProvider := NewCachedLLMProvider(mockProvider)
	ctx := context.Background()
	text := "Sample reference letter text"
	expectedData := &service.CredibilityData{
		Quotes:    []string{"Excellent work", "Great team player"},
		Sentiment: "POSITIVE",
	}

	mockProvider.On("ExtractCredibility", ctx, text).Return(expectedData, nil).Once()

	// Act - First call (cache miss)
	result1, err1 := cachedProvider.ExtractCredibility(ctx, text)

	// Assert - First call
	require.NoError(t, err1)
	assert.Equal(t, expectedData, result1)
	mockProvider.AssertExpectations(t)

	// Act - Second call (cache hit)
	result2, err2 := cachedProvider.ExtractCredibility(ctx, text)

	// Assert - Second call should return cached result
	require.NoError(t, err2)
	assert.Equal(t, expectedData, result2)
	assert.Equal(t, result1, result2)  // Same pointer means it's cached
	mockProvider.AssertExpectations(t) // Should not have been called again
}

func TestCachedLLMProvider_TailorProfile_WhenCacheMiss_ThenCallsProviderAndCaches(t *testing.T) {
	// Arrange
	mockProvider := new(servicemocks.MockLLMProvider)
	cachedProvider := NewCachedLLMProvider(mockProvider)
	ctx := context.Background()
	profileText := "Profile summary"
	jobDescription := "Job description"
	expectedSummary := "Tailored summary"
	expectedScore := 0.85

	mockProvider.On("TailorProfile", ctx, profileText, jobDescription).Return(expectedSummary, expectedScore, nil).Once()

	// Act - First call (cache miss)
	summary1, score1, err1 := cachedProvider.TailorProfile(ctx, profileText, jobDescription)

	// Assert - First call
	require.NoError(t, err1)
	assert.Equal(t, expectedSummary, summary1)
	assert.Equal(t, expectedScore, score1)
	mockProvider.AssertExpectations(t)

	// Act - Second call (cache hit)
	summary2, score2, err2 := cachedProvider.TailorProfile(ctx, profileText, jobDescription)

	// Assert - Second call should return cached result
	require.NoError(t, err2)
	assert.Equal(t, expectedSummary, summary2)
	assert.Equal(t, expectedScore, score2)
	mockProvider.AssertExpectations(t) // Should not have been called again
}

func TestCachedLLMProvider_ExtractProfileData_WhenDifferentTexts_ThenCachesSeparately(t *testing.T) {
	// Arrange
	mockProvider := new(servicemocks.MockLLMProvider)
	cachedProvider := NewCachedLLMProvider(mockProvider)
	ctx := context.Background()
	text1 := "First reference letter"
	text2 := "Second reference letter"
	expectedData1 := &service.ExtractedProfileData{
		CompanyName: "Company 1",
		Role:        "Role 1",
	}
	expectedData2 := &service.ExtractedProfileData{
		CompanyName: "Company 2",
		Role:        "Role 2",
	}

	mockProvider.On("ExtractProfileData", ctx, text1).Return(expectedData1, nil).Once()
	mockProvider.On("ExtractProfileData", ctx, text2).Return(expectedData2, nil).Once()

	// Act
	result1, err1 := cachedProvider.ExtractProfileData(ctx, text1)
	result2, err2 := cachedProvider.ExtractProfileData(ctx, text2)

	// Assert
	require.NoError(t, err1)
	require.NoError(t, err2)
	assert.Equal(t, expectedData1, result1)
	assert.Equal(t, expectedData2, result2)
	assert.NotEqual(t, result1, result2)

	// Act - Call again with same texts (should use cache)
	result1Cached, err1Cached := cachedProvider.ExtractProfileData(ctx, text1)
	result2Cached, err2Cached := cachedProvider.ExtractProfileData(ctx, text2)

	// Assert - Should return cached results
	require.NoError(t, err1Cached)
	require.NoError(t, err2Cached)
	assert.Equal(t, result1, result1Cached)
	assert.Equal(t, result2, result2Cached)
	mockProvider.AssertExpectations(t) // Should not have been called again
}

func TestCachedLLMProvider_TailorProfile_WhenDifferentInputs_ThenCachesSeparately(t *testing.T) {
	// Arrange
	mockProvider := new(servicemocks.MockLLMProvider)
	cachedProvider := NewCachedLLMProvider(mockProvider)
	ctx := context.Background()
	profileText1 := "Profile 1"
	profileText2 := "Profile 2"
	jobDesc1 := "Job 1"
	jobDesc2 := "Job 2"

	mockProvider.On("TailorProfile", ctx, profileText1, jobDesc1).Return("Summary 1", 0.8, nil).Once()
	mockProvider.On("TailorProfile", ctx, profileText2, jobDesc2).Return("Summary 2", 0.9, nil).Once()
	mockProvider.On("TailorProfile", ctx, profileText1, jobDesc2).Return("Summary 3", 0.7, nil).Once()

	// Act
	summary1, score1, err1 := cachedProvider.TailorProfile(ctx, profileText1, jobDesc1)
	summary2, score2, err2 := cachedProvider.TailorProfile(ctx, profileText2, jobDesc2)
	summary3, score3, err3 := cachedProvider.TailorProfile(ctx, profileText1, jobDesc2)

	// Assert
	require.NoError(t, err1)
	require.NoError(t, err2)
	require.NoError(t, err3)
	assert.Equal(t, "Summary 1", summary1)
	assert.Equal(t, 0.8, score1)
	assert.Equal(t, "Summary 2", summary2)
	assert.Equal(t, 0.9, score2)
	assert.Equal(t, "Summary 3", summary3)
	assert.Equal(t, 0.7, score3)

	// Act - Call again (should use cache)
	summary1Cached, score1Cached, err1Cached := cachedProvider.TailorProfile(ctx, profileText1, jobDesc1)

	// Assert - Should return cached result
	require.NoError(t, err1Cached)
	assert.Equal(t, "Summary 1", summary1Cached)
	assert.Equal(t, 0.8, score1Cached)
	mockProvider.AssertExpectations(t) // Should not have been called again
}

func TestCachedLLMProvider_ConcurrentAccess_WhenMultipleGoroutines_ThenThreadSafe(t *testing.T) {
	// Arrange
	mockProvider := new(servicemocks.MockLLMProvider)
	cachedProvider := NewCachedLLMProvider(mockProvider)
	ctx := context.Background()
	text := "Concurrent test text"
	expectedData := &service.ExtractedProfileData{
		CompanyName: "Test Company",
		Role:        "Engineer",
	}

	numGoroutines := 10
	// Allow 1 to numGoroutines calls (due to race conditions, multiple goroutines may hit cache miss simultaneously)
	// But all should return the same result
	mockProvider.On("ExtractProfileData", ctx, text).Return(expectedData, nil)

	// Act - Launch multiple goroutines
	var wg sync.WaitGroup
	results := make([]*service.ExtractedProfileData, numGoroutines)
	errors := make([]error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			results[idx], errors[idx] = cachedProvider.ExtractProfileData(ctx, text)
		}(i)
	}

	wg.Wait()

	// Assert - All should succeed and return consistent results
	// Note: Due to race conditions, multiple provider calls may occur, but results should be consistent
	for i := 0; i < numGoroutines; i++ {
		require.NoError(t, errors[i])
		assert.NotNil(t, results[i])
		// Verify all results have the same content (even if different pointers due to race)
		assert.Equal(t, expectedData.CompanyName, results[i].CompanyName)
		assert.Equal(t, expectedData.Role, results[i].Role)
	}

	// Verify provider was called at least once (cache should reduce calls, but race conditions may cause multiple)
	// The important thing is that all results are consistent
	calls := mockProvider.Calls
	assert.GreaterOrEqual(t, len(calls), 1, "Provider should be called at least once")
	assert.LessOrEqual(t, len(calls), numGoroutines, "Provider should not be called more than number of goroutines")
}

func TestCachedLLMProvider_hashText_WhenSameInput_ThenReturnsSameHash(t *testing.T) {
	// Arrange
	provider := NewCachedLLMProvider(nil)
	text := "test text"

	// Act
	hash1 := provider.hashText(text)
	hash2 := provider.hashText(text)

	// Assert
	assert.Equal(t, hash1, hash2)
	assert.NotEmpty(t, hash1)
}

func TestCachedLLMProvider_hashText_WhenDifferentInputs_ThenReturnsDifferentHashes(t *testing.T) {
	// Arrange
	provider := NewCachedLLMProvider(nil)
	text1 := "test text 1"
	text2 := "test text 2"

	// Act
	hash1 := provider.hashText(text1)
	hash2 := provider.hashText(text2)

	// Assert
	assert.NotEqual(t, hash1, hash2)
}

func TestCachedLLMProvider_hashTexts_WhenSameInputs_ThenReturnsSameHash(t *testing.T) {
	// Arrange
	provider := NewCachedLLMProvider(nil)
	text1 := "text 1"
	text2 := "text 2"

	// Act
	hash1 := provider.hashTexts(text1, text2)
	hash2 := provider.hashTexts(text1, text2)

	// Assert
	assert.Equal(t, hash1, hash2)
}

func TestCachedLLMProvider_hashTexts_WhenDifferentOrder_ThenReturnsDifferentHashes(t *testing.T) {
	// Arrange
	provider := NewCachedLLMProvider(nil)
	text1 := "text 1"
	text2 := "text 2"

	// Act
	hash1 := provider.hashTexts(text1, text2)
	hash2 := provider.hashTexts(text2, text1)

	// Assert
	assert.NotEqual(t, hash1, hash2) // Order matters
}
