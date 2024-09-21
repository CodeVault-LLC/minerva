package security

import (
	"net/http"
	"strings"
)

// ScanCors scans the Cross-Origin Resource Sharing headers of a given URL and returns a list of CORS headers
func ScanCors(header http.Header) []string {
	var corsIssues []string

	originHeader := header.Get("Access-Control-Allow-Origin")
	methodsHeader := header.Get("Access-Control-Allow-Methods")
	credentialsHeader := header.Get("Access-Control-Allow-Credentials")

	// Check if wildcard origin is allowed
	if originHeader == "*" {
		corsIssues = append(corsIssues, "Wildcard origin")
	}

	// Check if credentials are allowed
	if credentialsHeader == "true" && originHeader == "*" {
		corsIssues = append(corsIssues, "Wildcard origin and credentials")
	}

	// Check if all methods are allowed
	if strings.Contains(methodsHeader, "*") {
		corsIssues = append(corsIssues, "Wildcard methods")
	}

	return corsIssues
}
