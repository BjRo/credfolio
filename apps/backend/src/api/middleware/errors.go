package middleware

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"runtime/debug"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

// ErrorHandler recovers from panics and returns a JSON 500 error
func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				slog.Error("panic recovered", "error", rvr, "stack", string(debug.Stack()))

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(ErrorResponse{Error: "Internal Server Error"})
			}
		}()
		next.ServeHTTP(w, r)
	})
}
