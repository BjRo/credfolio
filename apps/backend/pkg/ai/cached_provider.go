package ai

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"sync"

	"github.com/credfolio/apps/backend/internal/service"
)

// cachedResult represents a cached LLM response
type cachedResult struct {
	profileData     *service.ExtractedProfileData
	credibilityData *service.CredibilityData
	tailorResult    *tailorResult
}

// tailorResult represents the result of TailorProfile
type tailorResult struct {
	summary    string
	matchScore float64
}

// CachedLLMProvider wraps an LLMProvider with caching
type CachedLLMProvider struct {
	provider service.LLMProvider
	cache    sync.Map // map[string]*cachedResult
}

// NewCachedLLMProvider creates a new cached LLM provider
func NewCachedLLMProvider(provider service.LLMProvider) *CachedLLMProvider {
	return &CachedLLMProvider{
		provider: provider,
	}
}

// hashText creates a SHA256 hash of the input text for use as a cache key
func (c *CachedLLMProvider) hashText(text string) string {
	hash := sha256.Sum256([]byte(text))
	return hex.EncodeToString(hash[:])
}

// hashTexts creates a combined hash of multiple text inputs
func (c *CachedLLMProvider) hashTexts(texts ...string) string {
	combined := ""
	for _, text := range texts {
		combined += text + "\n"
	}
	return c.hashText(combined)
}

// ExtractProfileData extracts structured profile data from text using cached results when available
func (c *CachedLLMProvider) ExtractProfileData(ctx context.Context, text string) (*service.ExtractedProfileData, error) {
	cacheKey := "profile_data:" + c.hashText(text)

	// Check cache
	if cached, ok := c.cache.Load(cacheKey); ok {
		if result, ok := cached.(*cachedResult); ok && result.profileData != nil {
			return result.profileData, nil
		}
	}

	// Call underlying provider
	data, err := c.provider.ExtractProfileData(ctx, text)
	if err != nil {
		return nil, err
	}

	// Cache the result
	c.cache.Store(cacheKey, &cachedResult{profileData: data})
	return data, nil
}

// ExtractCredibility extracts positive quotes and sentiment from text using cached results when available
func (c *CachedLLMProvider) ExtractCredibility(ctx context.Context, text string) (*service.CredibilityData, error) {
	cacheKey := "credibility:" + c.hashText(text)

	// Check cache
	if cached, ok := c.cache.Load(cacheKey); ok {
		if result, ok := cached.(*cachedResult); ok && result.credibilityData != nil {
			return result.credibilityData, nil
		}
	}

	// Call underlying provider
	data, err := c.provider.ExtractCredibility(ctx, text)
	if err != nil {
		return nil, err
	}

	// Cache the result
	c.cache.Store(cacheKey, &cachedResult{credibilityData: data})
	return data, nil
}

// TailorProfile generates a tailored profile summary based on job description using cached results when available
func (c *CachedLLMProvider) TailorProfile(ctx context.Context, profileText string, jobDescription string) (string, float64, error) {
	cacheKey := "tailor:" + c.hashTexts(profileText, jobDescription)

	// Check cache
	if cached, ok := c.cache.Load(cacheKey); ok {
		if result, ok := cached.(*cachedResult); ok && result.tailorResult != nil {
			return result.tailorResult.summary, result.tailorResult.matchScore, nil
		}
	}

	// Call underlying provider
	summary, matchScore, err := c.provider.TailorProfile(ctx, profileText, jobDescription)
	if err != nil {
		return "", 0, err
	}

	// Cache the result
	c.cache.Store(cacheKey, &cachedResult{
		tailorResult: &tailorResult{
			summary:    summary,
			matchScore: matchScore,
		},
	})
	return summary, matchScore, nil
}
