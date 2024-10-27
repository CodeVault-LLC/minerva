package responder

import (
	"encoding/json"
	"net/http"
)

type ResponseType string

const (
	ResponseTypeSuccess ResponseType = "success"
	ResponseTypeError   ResponseType = "error"
)

type APIResponse struct {
	Type       ResponseType `json:"type"`
	StatusCode int          `json:"status_code"`
	Message    string       `json:"message"`
	Data       interface{}  `json:"data,omitempty"`  // Populated in case of success responses.
	Error      *APIError    `json:"error,omitempty"` // Populated in case of error responses.
}

// APIError represents a structured error message.
type APIError struct {
	Code        string `json:"code"`        // Error code string (e.g., "auth_invalid_token").
	Description string `json:"description"` // User-friendly description of the error.
	Hint        string `json:"hint"`        // Optional hint for the user on how to resolve the issue.
}

// CreateSuccessResponse generates a success response with data.
func CreateSuccessResponse(data interface{}, message string) APIResponse {
	return APIResponse{
		Type:       ResponseTypeSuccess,
		StatusCode: http.StatusOK,
		Message:    message,
		Data:       data,
	}
}

// CreateErrorResponse generates an error response with a specified status code.
func CreateErrorResponse(code string, description string, hint string, statusCode int) APIResponse {
	return APIResponse{
		Type:       ResponseTypeError,
		StatusCode: statusCode,
		Message:    "An error occurred",
		Error: &APIError{
			Code:        code,
			Description: description,
			Hint:        hint,
		},
	}
}

// WriteJSONResponse writes the API response as JSON to the HTTP response writer.
func WriteJSONResponse(w http.ResponseWriter, apiResponse APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")


	w.WriteHeader(apiResponse.StatusCode)
	json.NewEncoder(w).Encode(apiResponse)
}
