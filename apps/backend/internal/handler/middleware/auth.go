package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type ContextKey string

const UserIDKey ContextKey = "userID"

// MockAuth is a mock authentication middleware that injects a hardcoded user ID
func MockAuth(next http.Handler) http.Handler {
	// Hardcoded user ID for development
	mockUserID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), UserIDKey, mockUserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserID extracts the user ID from the request context
func GetUserID(r *http.Request) uuid.UUID {
	userID, ok := r.Context().Value(UserIDKey).(uuid.UUID)
	if !ok {
		return uuid.Nil
	}
	return userID
}

// WithUserID returns a context with the user ID set (helper for tests)
func WithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}
