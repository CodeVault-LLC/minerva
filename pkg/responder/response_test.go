package responder_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/codevault-llc/minerva/pkg/responder"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestCreateSuccessResponse(t *testing.T) {
	data := map[string]string{"key": "value"}
	message := "Operation successful"
	response := responder.CreateSuccessResponse(data, message)

	assert.Equal(t, responder.ResponseTypeSuccess, response.Type)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, message, response.Message)
	assert.Equal(t, data, response.Data)
}

func TestWriteJSONResponse(t *testing.T) {
	app := fiber.New()

	app.Get("/test", func(c *fiber.Ctx) error {
		resp := responder.CreateSuccessResponse(map[string]string{"hello": "world"}, "success")
		responder.WriteJSONResponse(c, resp)
		return nil
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Check headers
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	assert.Equal(t, "nosniff", resp.Header.Get("X-Content-Type-Options"))
	assert.Equal(t, "DENY", resp.Header.Get("X-Frame-Options"))
	assert.Equal(t, "1; mode=block", resp.Header.Get("X-XSS-Protection"))

	// Check body
	var responseBody map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)
	assert.Equal(t, "success", responseBody["type"])
	assert.Equal(t, "success", responseBody["message"])
	assert.Equal(t, map[string]interface{}{"hello": "world"}, responseBody["data"])
}

func TestUnsupportedDataType(t *testing.T) {
	app := fiber.New()

	app.Get("/unsupported", func(c *fiber.Ctx) error {
		apiResponse := responder.APIResponse{
			Type:       responder.ResponseTypeSuccess,
			StatusCode: http.StatusOK,
			Message:    "Testing unsupported data type",
			Data:       struct{}{}, // Simulating unsupported type
		}
		responder.WriteJSONResponse(c, apiResponse)
		return nil
	})

	req := httptest.NewRequest(http.MethodGet, "/unsupported", nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	var body responder.APIResponse
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	assert.Equal(t, responder.ResponseTypeError, body.Type)
	assert.NotNil(t, body.Error)
	assert.Equal(t, "internal_server_error", body.Error.Code)
}

func TestErrorHandler(t *testing.T) {
	app := fiber.New(fiber.Config{
		ErrorHandler: responder.ErrorHandler,
	})

	app.Get("/test", func(c *fiber.Ctx) error {
		return &responder.APIError{
			Code:        "custom_error",
			Description: "This is a custom error",
			StatusCode:  http.StatusForbidden,
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, http.StatusForbidden, resp.StatusCode)

	// Check body
	var responseBody map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)
	assert.Equal(t, "error", responseBody["type"])
	assert.Equal(t, "This is a custom error", responseBody["error"].(map[string]interface{})["description"])
}
