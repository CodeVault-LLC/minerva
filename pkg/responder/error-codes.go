package responder

import (
	"net/http"
	"strconv"
)

// ErrorCode is a custom type for error codes.
type ErrorCode int

// List of error codes with clear naming conventions.
const (
	// Authentication Errors
	ErrAuthInvalidToken     ErrorCode = 40001
	ErrAuthPermissionDenied ErrorCode = 40002
	ErrAuthTokenExpired     ErrorCode = 40003
	ErrLimitReached         ErrorCode = 40005

	// Validation Errors
	ErrInvalidRequest         ErrorCode = 40004
	ErrValidationInvalidField ErrorCode = 40010
	ErrValidationMissingField ErrorCode = 40011
	ErrResourceNotFound       ErrorCode = 40400

	// Database Errors
	ErrDatabaseQueryFailed    ErrorCode = 50001
	ErrDatabaseConnection     ErrorCode = 50002
	ErrDatabaseRecordNotFound ErrorCode = 50003
	ErrScannerFailed          ErrorCode = 50004

	// Access & Generic Errors
	ErrForbiddenAccess     ErrorCode = 40300
	ErrInternalServerError ErrorCode = 50000
	ErrBadRequest          ErrorCode = 40000
)

// errorDetails holds descriptions and hints for each error code.
type errorDetails struct {
	Description string
	Hint        string
}

// Centralized error mappings with human-readable descriptions and hints.
var errorMapping = map[ErrorCode]errorDetails{
	ErrAuthInvalidToken:       {"Invalid token provided.", "Check if your token is valid and not expired."},
	ErrAuthPermissionDenied:   {"Permission denied for the requested resource.", "Ensure your account has the required permissions."},
	ErrAuthTokenExpired:       {"Authentication token has expired.", "Please login again to get a new token."},
	ErrValidationInvalidField: {"Invalid value in one or more fields.", "Check the input data for correctness."},
	ErrValidationMissingField: {"Required field is missing.", "Ensure that all mandatory fields are provided."},
	ErrDatabaseQueryFailed:    {"Database query execution failed.", "Try again later or contact support."},
	ErrDatabaseConnection:     {"Failed to connect to the database.", "Ensure the database is running and accessible."},
	ErrDatabaseRecordNotFound: {"No matching record found in the database.", "Check if the record exists before querying."},
	ErrForbiddenAccess:        {"Access to the requested resource is forbidden.", "Contact support if you believe this is a mistake."},
	ErrInternalServerError:    {"An internal server error occurred.", "Try again later or contact support."},
	ErrBadRequest:             {"The request could not be processed due to client error.", "Check the request format and parameters."},
	ErrInvalidRequest:         {"The request is invalid or malformed.", "Ensure the request is correctly formatted."},
	ErrResourceNotFound:       {"The requested resource was not found.", "Ensure the resource exists and the ID is correct."},
	ErrScannerFailed:          {"The scanner failed to process the request.", "Check the input data and try again."},
	ErrLimitReached:           {"Rate limit exceeded.", "Please wait before making more requests."},
}

// statusCodeMapping provides HTTP status codes based on error codes.
func statusCodeMapping(code ErrorCode) int {
	switch {
	case code == ErrForbiddenAccess:
		return http.StatusForbidden
	case code >= 40000 && code < 50000:
		return http.StatusBadRequest
	case code >= 50000 && code < 60000:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// NewAPIError creates an APIError with description and HTTP status code.
func NewAPIError(code ErrorCode) *APIError {
	details, exists := errorMapping[code]
	if !exists {
		details = errorDetails{
			Description: "An unknown error occurred.",
			Hint:        "Contact support with the provided error code.",
		}
	}
	return &APIError{
		Code:        strconv.Itoa(int(code)),
		Description: details.Description,
		Hint:        details.Hint,
		StatusCode:  statusCodeMapping(code),
	}
}

// CreateError generates an API response with error information.
func CreateError(code ErrorCode) APIResponse {
	apiErr := NewAPIError(code)
	return APIResponse{
		Type:       ResponseTypeError,
		StatusCode: apiErr.StatusCode,
		Message:    apiErr.Description,
		Error:      apiErr,
	}
}

// NewSuccessResponse generates a success API response with optional data.
func NewSuccessResponse(message string, data interface{}) APIResponse {
	return APIResponse{
		Type:       ResponseTypeSuccess,
		StatusCode: http.StatusOK,
		Message:    message,
		Data:       data,
	}
}

// Usage Example:
// errResponse := NewErrorResponse(ErrAuthInvalidToken)
// successResponse := NewSuccessResponse("Request successful", map[string]string{"key": "value"})
