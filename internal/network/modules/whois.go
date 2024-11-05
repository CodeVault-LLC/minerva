package modules

import (
	"github.com/codevault-llc/humblebrag-api/internal/core/models/entities"
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
	"github.com/likexian/whois"
	whoisparser "github.com/likexian/whois-parser"
)

type WhoisModule struct{}

func (m *WhoisModule) Run(job entities.JobModel) (interface{}, error) {
	whoisRecord, err := ScanWhois(utils.ConvertURLToDomain(job.URL))

	if err != nil || whoisRecord.Domain.Name == "" {
		return whoisparser.WhoisInfo{}, err
	} else {
		return whoisRecord, nil
	}
}

func (m *WhoisModule) Name() string {
	return "Whois"
}

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
