package helper

import (
	"testing"
)

func TestGenerateJWT(t *testing.T) {
	token, err := GenerateJWT(uint(1))

	if err != nil {
		t.Errorf("GenerateJWT() returned an error: %v", err)
	}

	if token == "" {
		t.Error("GenerateJWT() returned an empty token")
	}
}

func TestValidateJWT(t *testing.T) {
	token, err := GenerateJWT(uint(1))

	if err != nil {
		t.Errorf("GenerateJWT() returned an error: %v", err)
	}

	claims, err := ValidateJWT(token)

	if err != nil {
		t.Errorf("ValidateJWT() returned an error: %v", err)
	}

	if claims["id"] != float64(1) {
		t.Errorf("ValidateJWT() returned an invalid id: %v", claims["id"])
	}
}
