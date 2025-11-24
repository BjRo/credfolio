import type { ApiError } from "../api/generated/core/ApiError";

/**
 * Structured error response from backend
 */
interface ErrorResponse {
	error_id: number;
	message: string;
}

/**
 * Error ID to user-friendly message mapping
 * These messages are the single source of truth for user-facing error messages
 */
const ERROR_MESSAGES: Record<number, string> = {
	// Authentication & Authorization (1000-1099)
	1001: "Authentication required: Please log in and try again.",
	1002: "Access denied: You don't have permission to perform this action.",
	1003: "Invalid authentication token: Please log in again.",

	// Validation Errors (1100-1199)
	1101: "Invalid request: Please check your input and try again.",
	1102: "Invalid request body: The request format is incorrect.",
	1103: "Missing required field: Please fill in all required fields.",
	1104: "Invalid file type: Only PDF files are accepted.",
	1105: "File too large: Please upload a file smaller than 10MB.",
	1106: "Invalid job match ID: The provided job match ID is not valid.",
	1107: "Profile ID mismatch: The profile ID does not match.",

	// Resource Not Found (1200-1299)
	1201: "Profile not found: Please generate your profile first from the generate page.",
	1202: "Reference letter not found: The requested reference letter doesn't exist.",
	1203: "Job match not found: The requested job match doesn't exist.",
	1204: "Work experience not found: The requested work experience doesn't exist.",

	// Business Logic Errors (1300-1399)
	1301: "No reference letters found: Please upload at least one reference letter before generating your profile.",
	1302: "Job description is required: Please provide a job description to tailor your profile.",
	1303: "Job match mismatch: The job match does not belong to this profile.",

	// Processing Errors (1400-1499)
	1401: "Failed to process PDF: Unable to extract text from the uploaded file. Please ensure it's a valid PDF.",
	1402: "Failed to generate CV: Unable to create the PDF file. Please try again.",
	1403: "Profile generation failed: Unable to generate your profile. Please try again or contact support if the issue persists.",
	1404: "Profile tailoring failed: Unable to match your profile to the job description. Please try again.",
	1405: "Profile update failed: Unable to update your profile. Please try again.",

	// Server Errors (1500-1599)
	1501: "Server error: Something went wrong on our end. Please try again in a few moments.",
	1502: "Database error: Unable to access data. Please try again later.",
	1503: "External service error: A required service is temporarily unavailable. Please try again later.",
};

/**
 * Extracts error ID and message from ApiError body
 */
function extractStructuredError(error: ApiError): ErrorResponse | null {
	try {
		if (error.body && typeof error.body === "object") {
			// Check if body has the structured error format
			if ("error_id" in error.body && "message" in error.body) {
				return {
					error_id: error.body.error_id as number,
					message: error.body.message as string,
				};
			}
			// Try parsing if body is a string
			if (typeof error.body === "string") {
				const parsed = JSON.parse(error.body);
				if (parsed.error_id && parsed.message) {
					return {
						error_id: parsed.error_id,
						message: parsed.message,
					};
				}
			}
		}
	} catch {
		// If parsing fails, return null
	}
	return null;
}

/**
 * Extracts user-friendly error messages from API errors
 * Relies solely on structured error responses with error IDs from the backend
 */
export function getErrorMessage(error: unknown): string {
	// Check if it's an ApiError with structured error response
	if (error instanceof Error && "body" in error) {
		const apiError = error as ApiError;
		const structuredError = extractStructuredError(apiError);

		if (structuredError) {
			// Use error ID to get user-friendly message from mapping
			const mappedMessage = ERROR_MESSAGES[structuredError.error_id];
			if (mappedMessage) {
				return mappedMessage;
			}
			// Fall back to message from structured error if no mapping exists
			return structuredError.message;
		}
	}

	// Fallback for unknown error types (non-API errors, network errors, etc.)
	if (error instanceof Error) {
		return (
			error.message ||
			"An unexpected error occurred. Please try again or contact support if the issue persists."
		);
	}

	// Fallback for completely unknown error types
	return "An unexpected error occurred. Please try again or contact support if the issue persists.";
}
