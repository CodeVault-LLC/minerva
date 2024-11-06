package utils

import (
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

func IPsToStrings(ips []net.IP) []string {
	var result []string
	for _, ip := range ips {
		result = append(result, ip.String())
	}
	return result
}

func URIsToStrings(uris []*url.URL) []string {
	var result []string
	for _, uri := range uris {
		result = append(result, uri.String())
	}
	return result
}

func IPNetsToStrings(ipnets []*net.IPNet) []string {
	var result []string
	for _, ipnet := range ipnets {
		result = append(result, ipnet.String())
	}
	return result
}

func ConvertToStringSlice(rawMsg json.RawMessage) []string {
	var strSlice []string
	err := json.Unmarshal(rawMsg, &strSlice)
	if err != nil {
		fmt.Println("Failed to convert HTTPHeaders to []string:", err)
	}
	return strSlice
}

func StripProtocol(url string) string {
	if url[:8] == "https://" {
		return url[8:]
	} else if url[:7] == "http://" {
		return url[7:]
	}
	return url
}

func ConvertURLToDomain(inputURL string) string {
	u, err := url.Parse(inputURL)
	if err != nil {
		fmt.Println("Failed to parse URL:", err)
	}

	return u.Hostname()
}

// SafeString returns a default value if the string is empty
func SafeString(s string) string {
	if s == "" {
		return "N/A" // or any default value
	}
	return s
}

// IsNumeric checks if a string is a number
func IsNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// IsStylesheet checks if a link element is a stylesheet.
func IsStylesheet(node *html.Node) bool {
	for _, attr := range node.Attr {
		if attr.Key == "rel" && attr.Val == "stylesheet" {
			return true
		}
	}
	return false
}

// IsFont checks if a link element points to a font.
func IsFont(node *html.Node) bool {
	for _, attr := range node.Attr {
		if attr.Key == "rel" && strings.Contains(attr.Val, "font") {
			return true
		}
	}
	return false
}

func GenerateID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}

func GetCurrentTime() time.Time {
	return time.Now().UTC()
}
