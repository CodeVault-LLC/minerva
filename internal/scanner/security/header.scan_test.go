package security

import (
	"net/http"
	"testing"
)

func TestScanSecurityHeaders(t *testing.T) {
	var headers = http.Header{
		"Strict-Transport-Security": []string{"max-age=31536000"},
		"Content-Security-Policy":   []string{"default-src 'self'"},
		"X-Content-Type-Options":    []string{"nosniff"},
		"X-Frame-Options":           []string{"DENY"},
		"Referrer-Policy":           []string{"same-origin"},
		"Permissions-Policy":        []string{"geolocation=()"},
		"Cache-Control":             []string{"no-store"},
		"Pragma":                    []string{"no-cache"},
		"Feature-Policy":            []string{"geolocation 'self'"},
	}

	result := ScanSecurityHeaders(headers)

	// Test missing headers
	if len(result.MissingHeaders) != 0 {
		t.Errorf("Expected no missing headers, got %v", result.MissingHeaders)
	}

	// Test misconfigured headers
	if len(result.MisconfiguredHeaders) != 0 {
		t.Errorf("Expected no misconfigured headers, got %v", result.MisconfiguredHeaders)
	}
}
