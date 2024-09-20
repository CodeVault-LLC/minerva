package network

import (
	"github.com/likexian/whois"
	whoisparser "github.com/likexian/whois-parser"
)

// ScanWhois scans the whois information of a domain
func ScanWhois(domain string) (whoisparser.WhoisInfo, error) {
	result, err := whois.Whois(domain)
	if err != nil {
		return whoisparser.WhoisInfo{}, err
	}

	parsedResult, err := whoisparser.Parse(result)
	if err != nil {
		return whoisparser.WhoisInfo{}, err
	}

	return parsedResult, nil
}
