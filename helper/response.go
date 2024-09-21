package helper

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/codevault-llc/humblebrag-api/models"
	"github.com/codevault-llc/humblebrag-api/pkg/logger"
)

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err := w.Write(response)
	if err != nil {
		logger.Log.Error("Failed to write response: %v", err)
	}
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}

func DecodeJSON(body io.ReadCloser, v interface{}) {
	decoder := json.NewDecoder(body)
	err := decoder.Decode(v)
	if err != nil {
		logger.Log.Error("Failed to decode JSON: %v", err)
	}
}

func AddUserToContext(ctx context.Context, user models.UserModel) context.Context {
	return context.WithValue(ctx, "user", user)
}
