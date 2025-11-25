package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/credfolio/apps/backend/internal/handler/middleware"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMockAuthMiddleware_WhenRequestHasNoUser_InjectsMockUser(t *testing.T) {
	// Arrange
	var capturedUserID uuid.UUID
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.GetUserID(r.Context())
		if ok {
			capturedUserID = userID
		}
		w.WriteHeader(http.StatusOK)
	})

	wrapped := middleware.MockAuthMiddleware(handler)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	// Act
	wrapped.ServeHTTP(rr, req)

	// Assert
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.NotEqual(t, uuid.Nil, capturedUserID)
}

func TestMockAuthMiddleware_InjectsUserContext(t *testing.T) {
	// Arrange
	var capturedUser *middleware.UserContext
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := middleware.GetUser(r.Context())
		if ok {
			capturedUser = user
		}
		w.WriteHeader(http.StatusOK)
	})

	wrapped := middleware.MockAuthMiddleware(handler)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	// Act
	wrapped.ServeHTTP(rr, req)

	// Assert
	assert.NotNil(t, capturedUser)
	assert.Equal(t, "Development User", capturedUser.Name)
	assert.Equal(t, "dev@credfolio.app", capturedUser.Email)
}

func TestGetUserID_WhenUserInContext_ReturnsUserID(t *testing.T) {
	// Arrange
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.GetUserID(r.Context())
		assert.True(t, ok)
		assert.NotEqual(t, uuid.Nil, userID)
		w.WriteHeader(http.StatusOK)
	})

	wrapped := middleware.MockAuthMiddleware(handler)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	// Act
	wrapped.ServeHTTP(rr, req)

	// Assert
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestGetUserID_WhenNoUserInContext_ReturnsFalse(t *testing.T) {
	// Arrange
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	// Act
	userID, ok := middleware.GetUserID(req.Context())

	// Assert
	assert.False(t, ok)
	assert.Equal(t, uuid.Nil, userID)
}

func TestRequireAuth_WhenUserPresent_CallsNextHandler(t *testing.T) {
	// Arrange
	handlerCalled := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	})

	// First apply mock auth, then require auth
	wrapped := middleware.MockAuthMiddleware(middleware.RequireAuth(handler))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	// Act
	wrapped.ServeHTTP(rr, req)

	// Assert
	assert.True(t, handlerCalled)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestRequireAuth_WhenNoUser_Returns401(t *testing.T) {
	// Arrange
	handlerCalled := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	})

	// Only require auth, no mock auth
	wrapped := middleware.RequireAuth(handler)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	// Act
	wrapped.ServeHTTP(rr, req)

	// Assert
	assert.False(t, handlerCalled)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}
