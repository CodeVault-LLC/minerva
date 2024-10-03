package utils

import (
	"fmt"
	"net/url"
	"strings"

	regexp "github.com/wasilibs/go-re2"
)

// NormalizeURL formats the URL into a standard form, adds 'https' if necessary, removes 'www.' and trailing slashes.
func NormalizeURL(input string) string {
	urlNormalized := strings.TrimSpace(input)

	// Add "https://" if no scheme is provided
	if !strings.HasPrefix(urlNormalized, "http://") && !strings.HasPrefix(urlNormalized, "https://") {
		urlNormalized = "https://" + urlNormalized
	}

	// Remove 'www.' if present
	urlNormalized = strings.Replace(urlNormalized, "www.", "", 1)
	urlNormalized = strings.Replace(urlNormalized, ":///", "://", 1)

	parsedURL, err := url.Parse(urlNormalized)
	if err != nil {
		fmt.Println("Error parsing URL during normalization:", err)
		return ""
	}

	// Rebuild URL to remove any extraneous parts like paths or parameters
	urlNormalized = fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)

	// Remove trailing slash
	urlNormalized = strings.TrimSuffix(urlNormalized, "/")

	return urlNormalized
}

// ValidateURL checks if a URL is valid, looking for malformed URLs, SQL injection patterns, and illegal characters.
func ValidateURL(input string) bool {
	// Trim spaces and ensure it's not empty
	urlTrimmed := strings.TrimSpace(input)
	if urlTrimmed == "" {
		return false
	}

	// Check for SQL injection patterns or illegal characters
	sqlInjectionPattern := regexp.MustCompile(`(?i)('|\b(SELECT|INSERT|UPDATE|DELETE|DROP|UNION|ALTER|CREATE|EXEC)\b)`)
	illegalCharsPattern := regexp.MustCompile(`[<>]`)
	if sqlInjectionPattern.MatchString(urlTrimmed) || illegalCharsPattern.MatchString(urlTrimmed) {
		fmt.Println("Invalid characters or SQL injection attempt detected.")
		return false
	}

	// Parse the URL to ensure it's well-formed
	parsedURL, err := url.ParseRequestURI(urlTrimmed)
	if err != nil {
		fmt.Println("URL is malformed:", err)
		return false
	}

	// Ensure the scheme is either http or https
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		fmt.Println("URL scheme is invalid, only 'http' and 'https' are allowed.")
		return false
	}

	// Additional validation for domain name
	if parsedURL.Host == "" {
		fmt.Println("URL host is missing.")
		return false
	}

	fmt.Println("Valid URL:", parsedURL.String())
	return true
}
