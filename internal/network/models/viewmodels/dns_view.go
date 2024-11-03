package viewmodels

import "github.com/codevault-llc/humblebrag-api/internal/network/models/entities"

type DNS struct {
	ID uint `json:"id"`

	CNAME       []string `json:"cname"`
	ARecords    []string `json:"a_records"`
	AAAARecords []string `json:"aaaa_records"`
	MXRecords   []string `json:"mx_records"`
	NSRecords   []string `json:"ns_records"`
	TXTRecords  []string `json:"txt_records"`
	PTRRecord   string   `json:"ptr_record"`
	DNSSEC      bool     `json:"dnssec"`
}

func ConvertDNS(dns entities.DNSModel) DNS {
	return DNS{
		ID: dns.ID,

		CNAME:       dns.CNAME,
		ARecords:    dns.ARecords,
		AAAARecords: dns.AAAARecords,
		MXRecords:   dns.MXRecords,
		NSRecords:   dns.NSRecords,
		TXTRecords:  dns.TXTRecords,
		PTRRecord:   dns.PTRRecord,
		DNSSEC:      dns.DNSSEC,
	}
}
