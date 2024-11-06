package lists

import (
	"net"
)

// Check if the IP is a localhost or reserved domain
func isLocalhost(s string) bool {
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

/*
// Check if the IP is a multicast address (224.0.0.0/4)
func isMulticastIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil && parsedIP.IsMulticast()
}

// Check if the IP is a bogon (non-routable/reserved)
func isBogonIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	// List of common bogon ranges
	bogonRanges := []net.IPNet{
		{IP: net.IPv4(0, 0, 0, 0), Mask: net.CIDRMask(8, 32)},
		{IP: net.IPv4(100, 64, 0, 0), Mask: net.CIDRMask(10, 32)},   // Shared Address Space
		{IP: net.IPv4(192, 0, 0, 0), Mask: net.CIDRMask(24, 32)},    // IETF Protocol Assignments
		{IP: net.IPv4(198, 18, 0, 0), Mask: net.CIDRMask(15, 32)},   // Benchmarking
		{IP: net.IPv4(198, 51, 100, 0), Mask: net.CIDRMask(24, 32)}, // Documentation (TEST-NET-2)
	}

	for _, r := range bogonRanges {
		if r.Contains(parsedIP) {
			return true
		}
	}
	return false
}

// Check if the IP is a public IP
func isPublicIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	// Exclude private, loopback, multicast, or bogon IPs
	return !isPrivateIP(ip) && !parsedIP.IsLoopback() && !parsedIP.IsMulticast() && !isBogonIP(ip)
}

// Example to check if an IP is IPv6
func isIPv6(ip string) bool {
	return strings.Contains(ip, ":")
}
*/
