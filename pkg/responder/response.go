package responder

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

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
func createErrorResponse(code string, description string, statusCode int) APIResponse {
	return APIResponse{
		Type:       ResponseTypeError,
		StatusCode: statusCode,
		Message:    "An error occurred",
		Error: &APIError{
			Code:        code,
			Description: description,
		},
	}
}

// WriteJSONResponse writes the API response as JSON to the HTTP response writer.
func WriteJSONResponse(c *fiber.Ctx, apiResponse APIResponse) {
	if apiResponse.Data != nil {
		dataType := reflect.TypeOf(apiResponse.Data)

		if strings.Contains(dataType.String(), "entities") {
			logger.Log.Error("Failed to write JSON response", zap.Error(fmt.Errorf("Data type is not supported")))
			WriteJSONResponse(c, createErrorResponse("internal_server_error", "An internal server error occurred.", http.StatusInternalServerError))
			return
		}
	}

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
			WriteJSONResponse(c, createErrorResponse("internal_server_error", "An internal server error occurred.", http.StatusInternalServerError))
			return err
		}

		if apiError.Code == fiber.ErrNotFound.Code {
			WriteJSONResponse(c, createErrorResponse("not_found", "The requested resource was not found.", http.StatusNotFound))
			return nil
		}

		WriteJSONResponse(c, createErrorResponse("internal_server_error", "An internal server error occurred.", http.StatusInternalServerError))
		return nil
	}

	apiResponse := createErrorResponse(apiError.Code, apiError.Description, apiError.StatusCode)
	WriteJSONResponse(c, apiResponse)
	return nil
}
