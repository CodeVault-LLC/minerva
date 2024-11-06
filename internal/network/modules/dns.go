package modules

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/codevault-llc/minerva/internal/core/models/entities"
	"github.com/miekg/dns"
)

// GetCNAME retrieves CNAME records associated with the given URL.
func GetCNAME(url string) ([]string, error) {
	cname, err := net.LookupCNAME(url)
	if err != nil {
		return nil, err
	}

	return []string{cname}, nil
}

// GetARecords retrieves IPv4 A records for the given URL.
func GetARecords(url string) ([]string, error) {
	ips, err := net.LookupIP(url)
	if err != nil {
		return nil, err
	}

	var aRecords []string
	for _, ip := range ips {
		if ip.To4() != nil { // Filter for IPv4 addresses only
			aRecords = append(aRecords, ip.String())
		}
	}

	return aRecords, nil
}

// GetAAAARecords retrieves IPv6 AAAA records for the given URL.
func GetAAAARecords(url string) ([]string, error) {
	ips, err := net.LookupIP(url)
	if err != nil {
		return nil, err
	}

	var aaaaRecords []string
	for _, ip := range ips {
		if ip.To16() != nil && ip.To4() == nil { // Filter for IPv6 addresses only
			aaaaRecords = append(aaaaRecords, ip.String())
		}
	}

	return aaaaRecords, nil
}

// GetMXRecords retrieves MX (Mail Exchange) records for the given URL.
func GetMXRecords(url string) ([]string, error) {
	mxRecords, err := net.LookupMX(url)
	if err != nil {
		return nil, err
	}

	var mxList []string
	for _, mx := range mxRecords {
		mxList = append(mxList, fmt.Sprintf("%s (priority %d)", mx.Host, mx.Pref))
	}

	return mxList, nil
}

// GetNSRecords retrieves NS (Name Server) records for the given URL.
func GetNSRecords(url string) ([]string, error) {
	nsRecords, err := net.LookupNS(url)
	if err != nil {
		return nil, err
	}

	var nsList []string
	for _, ns := range nsRecords {
		nsList = append(nsList, ns.Host)
	}

	return nsList, nil
}

// GetTXTRecords retrieves TXT records (SPF, DKIM, DMARC, etc.).
func GetTXTRecords(url string) ([]string, error) {
	txtRecords, err := net.LookupTXT(url)
	if err != nil {
		return nil, err
	}

	return txtRecords, nil
}

// GetPTRRecord retrieves PTR (reverse DNS) record for a given IP.
func GetPTRRecord(ip string) (string, error) {
	names, err := net.LookupAddr(ip)
	if err != nil {
		return "", err
	}

	if len(names) > 0 {
		return names[0], nil
	}
	return "", nil
}

// DNSSECCheck checks DNSSEC status for the given domain.
func DNSSECCheck(domain string) (bool, error) {
	c := new(dns.Client)
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), dns.TypeDNSKEY)

	// Query the authoritative nameservers for DNSKEY
	nsRecords, err := net.LookupNS(domain)
	if err != nil {
		return false, err
	}

	for _, ns := range nsRecords {
		r, _, err := c.Exchange(m, net.JoinHostPort(ns.Host, "53"))
		if err != nil {
			continue
		}

		for _, rr := range r.Answer {
			if dnskey, ok := rr.(*dns.DNSKEY); ok && dnskey.Flags == 257 { // 257 means ZSK/KSK
				return true, nil
			}
		}
	}

	return false, fmt.Errorf("DNSSEC not found or not supported")
}

// DNSResults holds comprehensive DNS scan results.
type DNSResults struct {
	CNAME       []string
	ARecords    []string
	AAAARecords []string
	MXRecords   []string
	NSRecords   []string
	TXTRecords  []string
	PTRRecord   string
	DNSSEC      bool
}

type DNSModule struct{}

func (m *DNSModule) Run(job entities.JobModel) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	results := DNSResults{}
	var err error

	// Strip protocol (http/https) from URL before querying DNS
	domain := strings.TrimPrefix(strings.TrimPrefix(job.URL, "http://"), "https://")

	// Parallel DNS queries
	done := make(chan bool)

	go func() {
		results.CNAME, _ = GetCNAME(domain)
		done <- true
	}()

	go func() {
		results.ARecords, _ = GetARecords(domain)
		done <- true
	}()

	go func() {
		results.AAAARecords, _ = GetAAAARecords(domain)
		done <- true
	}()

	go func() {
		results.MXRecords, _ = GetMXRecords(domain)
		done <- true
	}()

	go func() {
		results.NSRecords, _ = GetNSRecords(domain)
		done <- true
	}()

	go func() {
		results.TXTRecords, _ = GetTXTRecords(domain)
		done <- true
	}()

	go func() {
		// Checking PTR for the first A record if exists
		if len(results.ARecords) > 0 {
			results.PTRRecord, _ = GetPTRRecord(results.ARecords[0])
		}
		done <- true
	}()

	go func() {
		results.DNSSEC, _ = DNSSECCheck(domain)
		done <- true
	}()

	// Wait for all goroutines to finish
	for i := 0; i < 8; i++ {
		select {
		case <-done:
		case <-ctx.Done():
			return DNSResults{}, fmt.Errorf("DNS scan timed out")
		}
	}

	return results, err
}

func (m *DNSModule) Name() string {
	return "DNS"
}
