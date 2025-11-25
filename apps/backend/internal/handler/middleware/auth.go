package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

// ContextKey is a type for context keys used in middleware
type ContextKey string

const (
	// UserIDKey is the context key for the current user ID
	UserIDKey ContextKey = "user_id"
	// UserKey is the context key for the current user
	UserKey ContextKey = "user"
)

// UserContext holds user information extracted from the request
type UserContext struct {
	ID    uuid.UUID
	Email string
	Name  string
}

// MockAuthMiddleware injects a mock user for development/testing
// In production, this would be replaced with real authentication
func MockAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a mock user for development
		mockUser := &UserContext{
			ID:    uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			Email: "dev@credfolio.app",
			Name:  "Development User",
		}

		// Add user to context
		ctx := context.WithValue(r.Context(), UserIDKey, mockUser.ID)
		ctx = context.WithValue(ctx, UserKey, mockUser)

		// Call next handler with updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserID extracts the user ID from the request context
func GetUserID(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(UserIDKey).(uuid.UUID)
	return userID, ok
}

// GetUser extracts the full user context from the request context
func GetUser(ctx context.Context) (*UserContext, bool) {
	user, ok := ctx.Value(UserKey).(*UserContext)
	return user, ok
}

// RequireAuth middleware ensures a user is authenticated
// Returns 401 if no user is found in context
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := GetUserID(r.Context())
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
