package utils

import (
	"encoding/json"
	"net"
	"net/url"
	"testing"
)

func TestIPsToStrings(t *testing.T) {
	ips := []net.IP{
		net.ParseIP("127.0.0.1"),
		net.ParseIP("192.168.0.1"),
	}

	strings := IPsToStrings(ips)

	if len(strings) != 2 {
		t.Error("Expected 2 strings")
	}

	if strings[0] != "127.0.0.1" {
		t.Error("Expected '127.0.0.1'")
	}

	if strings[1] != "192.168.0.1" {
		t.Error("Expected '192.168.0.1'")
	}
}

func TestURIsToStrings(t *testing.T) {
	uris := []*url.URL{
		{Scheme: "http", Host: "example.com"},
		{Scheme: "https", Host: "example.org"},
	}

	strings := URIsToStrings(uris)

	if len(strings) != 2 {
		t.Error("Expected 2 strings")
	}

	if strings[0] != "http://example.com" {
		t.Error("Expected 'http://example.com'")
	}

	if strings[1] != "https://example.org" {
		t.Error("Expected 'https://example.org'")
	}
}
func TestIPNetsToStrings(t *testing.T) {
	ipnets := []*net.IPNet{
		{IP: net.ParseIP("192.168.0.0"), Mask: net.CIDRMask(24, 32)},
		{IP: net.ParseIP("10.0.0.0"), Mask: net.CIDRMask(16, 32)},
	}

	strings := IPNetsToStrings(ipnets)

	if len(strings) != 2 {
		t.Error("Expected 2 strings")
	}

	if strings[0] != "192.168.0.0/24" {
		t.Error("Expected '192.168.0.0/24'")
	}

	if strings[1] != "10.0.0.0/16" {
		t.Error("Expected '10.0.0.0/16'")
	}
}

func TestConvertToStringSlice(t *testing.T) {
	rawMsg := json.RawMessage(`["hello", "world"]`)

	strSlice := ConvertToStringSlice(rawMsg)

	if len(strSlice) != 2 {
		t.Error("Expected 2 strings")
	}

	if strSlice[0] != "hello" {
		t.Error("Expected 'hello'")
	}

	if strSlice[1] != "world" {
		t.Error("Expected 'world'")
	}
}

func TestStripProtocol(t *testing.T) {
	url1 := "https://example.com"
	stripped1 := StripProtocol(url1)
	if stripped1 != "example.com" {
		t.Errorf("Expected 'example.com', got '%s'", stripped1)
	}

	url2 := "http://www.example.org"
	stripped2 := StripProtocol(url2)
	if stripped2 != "www.example.org" {
		t.Errorf("Expected 'www.example.org', got '%s'", stripped2)
	}

	url3 := "ftp://ftp.example.net"
	stripped3 := StripProtocol(url3)
	if stripped3 != "ftp://ftp.example.net" {
		t.Errorf("Expected 'ftp://ftp.example.net', got '%s'", stripped3)
	}
}

func TestConvertURLToDomain(t *testing.T) {
	inputURL := "https://www.example.com"
	expectedDomain := "www.example.com"

	domain := ConvertURLToDomain(inputURL)

	if domain != expectedDomain {
		t.Errorf("Expected '%s', got '%s'", expectedDomain, domain)
	}

	inputURL = "http://subdomain.example.org"
	expectedDomain = "subdomain.example.org"

	domain = ConvertURLToDomain(inputURL)

	if domain != expectedDomain {
		t.Errorf("Expected '%s', got '%s'", expectedDomain, domain)
	}

	inputURL = "ftp://ftp.example.net"
	expectedDomain = "ftp.example.net"

	domain = ConvertURLToDomain(inputURL)

	if domain != expectedDomain {
		t.Errorf("Expected '%s', got '%s'", expectedDomain, domain)
	}
}

func TestSafeString(t *testing.T) {
	s := "hello"
	safe := SafeString(s)
	if safe != s {
		t.Errorf("Expected '%s', got '%s'", s, safe)
	}

	s = ""
	safe = SafeString(s)
	if safe != "N/A" {
		t.Errorf("Expected 'N/A', got '%s'", safe)
	}
}

// TestIsNumeric tests the IsNumeric function
func TestIsNumeric(t *testing.T) {
	// Test with a number
	if !IsNumeric("123") {
		t.Error("Expected '123' to be numeric")
	}

	// Test with a non-number
	if IsNumeric("abc") {
		t.Error("Expected 'abc' to not be numeric")
	}
}
