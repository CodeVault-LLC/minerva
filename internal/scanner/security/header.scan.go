package security

import (
	"fmt"
	"net/http"
	"strings"
)

// List of important security headers to check
var securityHeaders = []string{
	"Strict-Transport-Security", // HSTS header
	"Content-Security-Policy",   // CSP header
	"X-Content-Type-Options",    // Prevent MIME type sniffing
	"X-Frame-Options",           // Prevent clickjacking
	"Referrer-Policy",           // Control the referrer information
	"Permissions-Policy",        // Control which browser features a page can use
	"Cache-Control",             // Control caching
	"Pragma",                    // Cache control for older browsers
	"Feature-Policy",            // Deprecated, replaced by Permissions-Policy
}

type ScanSecurityHead struct {
	MissingHeaders       []string
	MisconfiguredHeaders []string
}

// ScanSecurityHeaders scans the provided HTTP headers and returns a ScanSecurityHead struct
// containing the missing and misconfigured security headers.
//
// Parameters:
//   - headers: The HTTP headers to be scanned.
//
// Returns:
//   - ScanSecurityHead: A struct containing the missing and misconfigured security headers.
//
// Example:
//
//	headers := http.Header{
//	  "Strict-Transport-Security": "max-age=31536000",
//	  "X-XSS-Protection":         "1; mode=block",
//	}
//	result := ScanSecurityHeaders(headers)
//
//	fmt.Println(result.MissingHeaders)             // Output: []
//	fmt.Println(result.MisconfiguredHeaders)       // Output: []
//
// Note:
//
//	The function checks for missing headers by comparing them against a predefined list of security headers.
//	It also checks for misconfigured headers such as "Strict-Transport-Security" without the "max-age" directive
//	and "X-XSS-Protection" with a value of "0".
func ScanSecurityHeaders(headers http.Header) ScanSecurityHead {
	var missingHeaders []string
	var misconfiguredHeaders []string

	for _, header := range securityHeaders {
		value := headers.Get(header)
		if value == "" {
			fmt.Println("Missing header: ", header, value)
			missingHeaders = append(missingHeaders, header)
		}
	}

	if hsts := headers.Get("Strict-Transport-Security"); hsts != "" && !strings.Contains(hsts, "max-age=") {
		misconfiguredHeaders = append(misconfiguredHeaders, "Strict-Transport-Security")
	}

	if xssProtection := headers.Get("X-XSS-Protection"); xssProtection == "0" {
		fmt.Println("X-XSS-Protection: 0")
		misconfiguredHeaders = append(misconfiguredHeaders, "X-XSS-Protection")
	}

	return ScanSecurityHead{
		MissingHeaders:       missingHeaders,
		MisconfiguredHeaders: misconfiguredHeaders,
	}
}
