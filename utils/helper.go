package utils

import (
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
