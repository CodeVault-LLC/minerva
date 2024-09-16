package ip

import (
	"net"
	"strings"

	"github.com/codevault-llc/humblebrag-api/utils"
)

// GetDNSNames retrieves CNAME records associated with the given URL.
func GetDNSNames(url string) ([]string, error) {
	cname, err := net.LookupCNAME(url)
	if err != nil {
		return nil, err
	}

	// Return the single CNAME (could add more if needed)
	return []string{cname}, nil
}

// GetPermittedDNS retrieves TXT records (e.g., SPF) that specify permitted DNS entries.
func GetPermittedDNS(url string) ([]string, error) {
	txtRecords, err := net.LookupTXT(url)
	if err != nil {
		return nil, err
	}

	var permitted []string
	for _, record := range txtRecords {
		if strings.Contains(record, "v=spf1") {
			// Extract the permitted domains or IPs from the SPF record
			permitted = append(permitted, extractPermittedFromSPF(record))
		}
	}

	return permitted, nil
}

// Helper function to extract permitted domains or IPs from an SPF record
func extractPermittedFromSPF(spfRecord string) string {
	// This is a simple placeholder. You can extend the logic to properly parse SPF records.
	return strings.TrimPrefix(spfRecord, "v=spf1 ")
}

// GetExcludedDNS looks for DNS records that might indicate exclusions (e.g., DNSBL).
func GetExcludedDNS(url string) ([]string, error) {
	// Example implementation for DNSBL (DNS-based blackhole list), which isn't part of Go's standard library
	// Placeholder - could query known DNSBL services or fetch denial records from external lists.
	// For now, just return a static list or a query like we did for SPF.

	// Placeholder: Implement actual blacklist lookup later
	return []string{"example.excluded-dns.com"}, nil
}

type DNSResults struct {
	CNAME     []string
	Permitted []string
	Excluded  []string
}

func GetDNSScan(url string) (DNSResults, error) {
	url = utils.StripProtocol(url)

	cname, err := GetDNSNames(url)
	if err != nil {
		return DNSResults{}, err
	}

	permitted, err := GetPermittedDNS(url)
	if err != nil {
		return DNSResults{}, err
	}

	excluded, err := GetExcludedDNS(url)
	if err != nil {
		return DNSResults{}, err
	}

	return DNSResults{
		CNAME:     cname,
		Permitted: permitted,
		Excluded:  excluded,
	}, nil
}
