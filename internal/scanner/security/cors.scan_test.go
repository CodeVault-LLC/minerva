package security

import (
	"net/http"
	"testing"
)

func TestScanCors(t *testing.T) {
	header := http.Header{
		"Access-Control-Allow-Origin":      []string{"*"},
		"Access-Control-Allow-Methods":     []string{"GET, POST, PUT"},
		"Access-Control-Allow-Credentials": []string{"true"},
	}

	expected := []string{"Wildcard origin", "Wildcard origin and credentials"}

	result := ScanCors(header)

	if len(result) != len(expected) {
		t.Errorf("Expected %d CORS issues, but got %d", len(expected), len(result))
	}

	for i, issue := range expected {
		if result[i] != issue {
			t.Errorf("Expected CORS issue '%s', but got '%s'", issue, result[i])
		}
	}
}

func TestScanCorsNoIssues(t *testing.T) {
	header := http.Header{
		"Access-Control-Allow-Origin":      []string{"https://example.com"},
		"Access-Control-Allow-Methods":     []string{"GET, POST, PUT"},
		"Access-Control-Allow-Credentials": []string{"true"},
	}

	result := ScanCors(header)

	if len(result) > 0 {
		t.Errorf("Expected no CORS issues, but got %d", len(result))
	}
}
