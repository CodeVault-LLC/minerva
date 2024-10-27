package responder

import (
	"net/http"
	"strconv"
)

// ErrorCode is a custom type for error codes.
type ErrorCode int

// List of error codes.
const (
	// Authentication error codes.
	ErrAuthInvalidToken     ErrorCode = 40001
	ErrAuthPermissionDenied ErrorCode = 40002
	ErrAuthTokenExpired     ErrorCode = 40003

	// Validation error codes.
	ErrInvalidRequest         ErrorCode = 40004
	ErrValidationInvalidField ErrorCode = 40010
	ErrValidationMissingField ErrorCode = 40011

	// Database error codes.
	ErrDatabaseQueryFailed    ErrorCode = 50001
	ErrDatabaseConnection     ErrorCode = 50002
	ErrDatabaseRecordNotFound ErrorCode = 50003

	ErrScannerFailed ErrorCode = 50004

	// Forbidden error code.
	ErrForbiddenAccess ErrorCode = 40300

	// Generic error codes.
	ErrInternalServerError ErrorCode = 50000
	ErrBadRequest          ErrorCode = 40000
)

// errorDetails holds the error metadata including descriptions and hints.
type errorDetails struct {
	Description string
	Hint        string
}

// errorMapping holds the details for each error code.
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
}

// statusCodeMapping derives HTTP status codes from error codes.
func statusCodeMapping(code ErrorCode) int {
	switch {
	case code >= 20000 && code < 30000:
		return http.StatusOK
	case code >= 30000 && code < 40000:
		return http.StatusUnauthorized
	case code >= 40000 && code < 50000:
		return http.StatusBadRequest
	case code >= 50000 && code < 60000:
		return http.StatusInternalServerError
	case code == ErrForbiddenAccess:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

// CreateError generates an error response based on the error code.
func CreateError(code ErrorCode) APIResponse {
	details, exists := errorMapping[code]
	if !exists {
		details = errorDetails{
			Description: "An unknown error occurred.",
			Hint:        "Contact support with the provided error code.",
		}
	}

	return CreateErrorResponse(strconv.Itoa(int(code)), details.Description, details.Hint, statusCodeMapping(code))
}
