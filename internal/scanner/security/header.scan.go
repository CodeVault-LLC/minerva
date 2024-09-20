package security

import (
	"net/http"
	"strings"
)

// List of important security headers to check
var securityHeaders = []string{
	"Strict-Transport-Security", // HSTS header
	"Content-Security-Policy",   // CSP header
	"X-Content-Type-Options",    // Prevent MIME type sniffing
	"X-Frame-Options",           // Prevent clickjacking
	"X-XSS-Protection",          // XSS protection
	"Referrer-Policy",           // Control the referrer information
	"Permissions-Policy",        // Control which browser features a page can use
	"Expect-CT",                 // Certificate Transparency
	"Cache-Control",             // Control caching
	"Pragma",                    // Cache control for older browsers
	"Feature-Policy",            // Deprecated, replaced by Permissions-Policy
}

type ScanSecurityHead struct {
	MissingHeaders       []string
	MisconfiguredHeaders []string
}

func ScanSecurityHeaders(headers http.Header) ScanSecurityHead {
	var missingHeaders []string
	var misconfiguredHeaders []string

	for _, header := range securityHeaders {
		value := headers.Get(header)
		if value == "" {
			missingHeaders = append(missingHeaders, header)
		}
	}

	if hsts := headers.Get("Strict-Transport-Security"); hsts != "" && !strings.Contains(hsts, "max-age=") {
		misconfiguredHeaders = append(misconfiguredHeaders, "Strict-Transport-Security")
	}

	if xssProtection := headers.Get("X-XSS-Protection"); xssProtection == "0" {
		misconfiguredHeaders = append(misconfiguredHeaders, "X-XSS-Protection")
	}

	return ScanSecurityHead{
		MissingHeaders:       missingHeaders,
		MisconfiguredHeaders: misconfiguredHeaders,
	}
}
