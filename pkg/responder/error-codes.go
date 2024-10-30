package responder

import (
	"net/http"
	"strconv"
)

// Define custom error and HTTP status codes
type ErrorCode int

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

type errorDetails struct {
	Description string
	StatusCode  int
}

var errorMapping = map[ErrorCode]errorDetails{
	ErrAuthInvalidToken:       {"Invalid token provided.", http.StatusBadRequest},
	ErrAuthPermissionDenied:   {"Permission denied for the requested resource.", http.StatusForbidden},
	ErrAuthTokenExpired:       {"Authentication token has expired.", http.StatusUnauthorized},
	ErrValidationInvalidField: {"Invalid field value.", http.StatusBadRequest},
	ErrValidationMissingField: {"Required field is missing.", http.StatusBadRequest},
	ErrDatabaseQueryFailed:    {"Database query failed.", http.StatusInternalServerError},
	ErrDatabaseConnection:     {"Database connection failed.", http.StatusInternalServerError},
	ErrDatabaseRecordNotFound: {"Record not found in the database.", http.StatusNotFound},
	ErrForbiddenAccess:        {"Forbidden access.", http.StatusForbidden},
	ErrInternalServerError:    {"An internal server error occurred.", http.StatusInternalServerError},
	ErrBadRequest:             {"Bad request format.", http.StatusBadRequest},
	ErrInvalidRequest:         {"Invalid or malformed request.", http.StatusBadRequest},
	ErrResourceNotFound:       {"Resource not found.", http.StatusNotFound},
	ErrScannerFailed:          {"Request scanning failed.", http.StatusInternalServerError},
	ErrLimitReached:           {"Rate limit exceeded.", http.StatusTooManyRequests},
}

// NewAPIError creates an APIError based on the ErrorCode, defaulting to a generic error if not mapped
func NewAPIError(code ErrorCode) *APIError {
	details, exists := errorMapping[code]
	if !exists {
		details = errorDetails{
			Description: "An unknown error occurred.",
			StatusCode:  http.StatusInternalServerError,
		}
	}
	return &APIError{
		Code:        strconv.Itoa(int(code)),
		Description: details.Description,
		StatusCode:  details.StatusCode,
	}
}

// CreateError generates an error response
func CreateError(code ErrorCode) APIResponse {
	apiErr := NewAPIError(code)
	return APIResponse{
		Type:       "error",
		StatusCode: apiErr.StatusCode,
		Message:    apiErr.Description,
		Error:      apiErr,
	}
}
