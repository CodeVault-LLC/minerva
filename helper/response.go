package helper

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/codevault-llc/humblebrag-api/models"
)

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}

func DecodeJSON(body io.ReadCloser, v interface{}) {
	decoder := json.NewDecoder(body)
	decoder.Decode(v)
}

func AddUserToContext(ctx context.Context, user models.UserModel) context.Context {
	return context.WithValue(ctx, "user", user)
}
