package utils

import "testing"

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "URL with http scheme",
			input:    "http://example.com",
			expected: "http://example.com",
		},
		{
			name:     "URL with https scheme",
			input:    "https://example.com",
			expected: "https://example.com",
		},
		{
			name:     "URL without scheme",
			input:    "example.com",
			expected: "https://example.com",
		},
		{
			name:     "URL with www prefix",
			input:    "https://www.example.com",
			expected: "https://example.com",
		},
		{
			name:     "URL with trailing slash",
			input:    "https://example.com/",
			expected: "https://example.com",
		},
		{
			name:     "URL with multiple slashes",
			input:    "https:///example.com",
			expected: "https://example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NormalizeURL(tt.input); got != tt.expected {
				t.Errorf("NormalizeURL() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Valid URL with http scheme",
			input:    "http://example.com",
			expected: true,
		},
		{
			name:     "Valid URL with https scheme",
			input:    "https://example.com",
			expected: true,
		},
		{
			name:     "Invalid URL with ftp scheme",
			input:    "ftp://example.com",
			expected: false,
		},
		{
			name:     "Invalid URL with no scheme",
			input:    "example.com",
			expected: false,
		},
		{
			name:     "Invalid URL with invalid characters",
			input:    "https://example.com/<>",
			expected: false,
		},
		{
			name:     "Invalid URL with SQL injection pattern",
			input:    "https://example.com/SELECT",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateURL(tt.input); got != tt.expected {
				t.Errorf("ValidateURL() = %v, want %v", got, tt.expected)
			}
		})
	}
}
