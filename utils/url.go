package utils

import (
	"fmt"
	"net/url"
	"strings"
)

func NormalizeURL(input string) string {
	urlNormalized := strings.TrimSpace(input)

	if !strings.HasPrefix(urlNormalized, "http://") && !strings.HasPrefix(urlNormalized, "https://") {
		urlNormalized = "https://" + urlNormalized
	}

	urlNormalized = strings.Replace(urlNormalized, "www.", "", 1)
	urlNormalized = strings.Replace(urlNormalized, ":///", "://", 1)

	parsedURL, err := url.Parse(urlNormalized)
	if err != nil {
		return urlNormalized
	}

	// Construct a new URL with only the scheme and host
	urlNormalized = parsedURL.Scheme + "://" + parsedURL.Host

	// Remove the last / from the URL
	urlNormalized = strings.TrimSuffix(urlNormalized, "/")

	fmt.Println("Normalized URL:", urlNormalized)

	return urlNormalized
}

// ValidateURL checks if the URL is well-formed and has either http or https as the scheme.
func ValidateURL(input string) bool {
	// Parse the URL to check if it's well-formed
	parsedURL, err := url.ParseRequestURI(input)
	if err != nil {
		return false
	}

	// Ensure the scheme is either http or https
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return false
	}

	return true
}
