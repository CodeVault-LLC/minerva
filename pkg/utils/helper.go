package utils

import (
	"encoding/json"
	"fmt"
	"net"
	"net/url"
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

func SafeString(s string) string {
	if s == "" {
		return "N/A" // or any default value
	}
	return s
}
