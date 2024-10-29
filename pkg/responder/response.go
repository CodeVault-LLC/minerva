package responder

import (
	"net/http"

	"github.com/codevault-llc/humblebrag-api/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
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
	StatusCode  int    `json:"status_code"` // HTTP status code.
}

// Error returns the error message of the APIError.
func (e *APIError) Error() string {
	return e.Description
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

// createErrorResponse generates an error response with a specified status code.
func createErrorResponse(code string, description string, hint string, statusCode int) APIResponse {
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
func WriteJSONResponse(c *fiber.Ctx, apiResponse APIResponse) {
	c.Response().Header.Set("Content-Type", "application/json")
	c.Response().Header.Set("X-Content-Type-Options", "nosniff")
	c.Response().Header.Set("X-Frame-Options", "DENY")
	c.Response().Header.Set("X-XSS-Protection", "1; mode=block")

	c.Status(apiResponse.StatusCode)
	err := c.JSON(apiResponse)
	if err != nil {
		logger.Log.Error("Failed to write JSON response", zap.Error(err))
	}
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	apiError, ok := err.(*APIError)
	if !ok {
		apiError, ok := err.(*fiber.Error)
		if !ok {
			logger.Log.Error("An internal server error occurred", zap.Error(err))
			WriteJSONResponse(c, createErrorResponse("internal_server_error", "An internal server error occurred.", "Try again later or contact support.", http.StatusInternalServerError))
			return err
		}

		if apiError.Code == fiber.ErrNotFound.Code {
			WriteJSONResponse(c, createErrorResponse("not_found", "The requested resource was not found.", "Ensure the resource exists and the ID is correct.", http.StatusNotFound))
			return nil
		}

		WriteJSONResponse(c, createErrorResponse("internal_server_error", "An internal server error occurred.", "Try again later or contact support.", http.StatusInternalServerError))
		return nil
	}

	apiResponse := createErrorResponse(apiError.Code, apiError.Description, apiError.Hint, apiError.StatusCode)
	WriteJSONResponse(c, apiResponse)
	return nil
}
