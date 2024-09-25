package helper

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"reflect"

	"github.com/codevault-llc/humblebrag-api/internal/database/models"
	"github.com/codevault-llc/humblebrag-api/pkg/logger"
)

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	if payload != nil && reflect.TypeOf(payload).Kind() == reflect.Slice && reflect.ValueOf(payload).Len() == 0 {
		payload = []interface{}{}
	}

	response, err := json.Marshal(payload)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to marshal JSON payload")
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	_, err = w.Write(response)
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

func AddLicenseToContext(ctx context.Context, license models.LicenseModel) context.Context {
	return context.WithValue(ctx, "license", license)
}
