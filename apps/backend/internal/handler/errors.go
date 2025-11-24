package handler

import (
	"encoding/json"
	"net/http"
)

// Error codes - numeric IDs for different error types
const (
	// Authentication & Authorization (1000-1099)
	ErrorCodeUnauthorized = 1001
	ErrorCodeForbidden    = 1002
	ErrorCodeInvalidToken = 1003

	// Validation Errors (1100-1199)
	ErrorCodeInvalidRequest     = 1101
	ErrorCodeInvalidRequestBody = 1102
	ErrorCodeMissingField       = 1103
	ErrorCodeInvalidFileType    = 1104
	ErrorCodeFileTooLarge       = 1105
	ErrorCodeInvalidJobMatchID  = 1106
	ErrorCodeProfileIDMismatch  = 1107

	// Resource Not Found (1200-1299)
	ErrorCodeProfileNotFound         = 1201
	ErrorCodeReferenceLetterNotFound = 1202
	ErrorCodeJobMatchNotFound        = 1203
	ErrorCodeWorkExperienceNotFound  = 1204

	// Business Logic Errors (1300-1399)
	ErrorCodeNoReferenceLetters     = 1301
	ErrorCodeJobDescriptionRequired = 1302
	ErrorCodeJobMatchMismatch       = 1303

	// Processing Errors (1400-1499)
	ErrorCodePDFExtractionFailed     = 1401
	ErrorCodePDFGenerationFailed     = 1402
	ErrorCodeProfileGenerationFailed = 1403
	ErrorCodeProfileTailoringFailed  = 1404
	ErrorCodeProfileUpdateFailed     = 1405

	// Server Errors (1500-1599)
	ErrorCodeInternalServerError  = 1501
	ErrorCodeDatabaseError        = 1502
	ErrorCodeExternalServiceError = 1503
)

// ErrorResponse represents a structured error response
type ErrorResponse struct {
	ErrorID int    `json:"error_id"`
	Message string `json:"message"`
}

// writeErrorResponse writes a structured JSON error response
func writeErrorResponse(w http.ResponseWriter, statusCode int, errorID int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := ErrorResponse{
		ErrorID: errorID,
		Message: message,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Fallback to plain text if JSON encoding fails
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Internal server error"))
	}
}
