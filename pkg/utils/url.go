package utils

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"unsafe"

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

// ParseUint safely parses a string to uint, with an upper bound check for 32-bit systems.
func ParseUint(input string) (uint, error) {
	parsed, err := strconv.ParseUint(input, 10, 64)
	if err != nil {
		return 0, err
	}

	// Check for overflow if `uint` is 32 bits (4 bytes)
	if unsafe.Sizeof(uint(0)) == 4 && parsed > uint64(^uint32(0)) {
		return 0, errors.New("value exceeds the maximum size of uint on a 32-bit system")
	}

	return uint(parsed), nil
}

// Check if the URL is local or remote. If local then return true, otherwise false.
func IsLocalURL(input string) bool {
	parsedURL, err := url.ParseRequestURI(input)
	if err != nil {
		return false
	}

	if parsedURL.Scheme == "file" || parsedURL.Scheme == "localhost" {
		return true
	}

	if isPrivateIP(net.ParseIP(parsedURL.Host)) {
		return true
	}

	if isLocalhost(parsedURL.Host) {
		return true
	}

	return false
}

// isPrivateIP checks if an IP address is private
func isPrivateIP(ip net.IP) bool {
	privateIPBlocks := []*net.IPNet{
		// IPv4 private addresses
		{IP: net.IPv4(10, 0, 0, 0), Mask: net.CIDRMask(8, 32)},
		{IP: net.IPv4(172, 16, 0, 0), Mask: net.CIDRMask(12, 32)},
		{IP: net.IPv4(192, 168, 0, 0), Mask: net.CIDRMask(16, 32)},
		// IPv6 private addresses
		{IP: net.ParseIP("fc00::"), Mask: net.CIDRMask(7, 128)},
	}

	for _, block := range privateIPBlocks {
		if block.Contains(ip) {
			return true
		}
	}

	return false
}

// isLocalhost checks if a string is a localhost or reserved domain
func isLocalhost(s string) bool {
	s = strings.TrimPrefix(s, "http://")
	s = strings.TrimPrefix(s, "https://")

	if strings.Contains(s, ":") {
		s = strings.Split(s, ":")[0]
	}

	if s == "localhost" || s == "localhost.localdomain" || s == "local" || s == "broadcasthost" {
		return true
	}

	if ip := net.ParseIP(s); ip != nil {
		if ip.IsLoopback() || ip.IsLinkLocalMulticast() || ip.IsLinkLocalUnicast() {
			return true
		}
	}

	return false
}
