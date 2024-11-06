package viewmodels

import "github.com/codevault-llc/minerva/internal/network/models/entities"

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

func ConvertDNS(dns entities.DnsModel) DNS {
	return DNS{
		ID: dns.Id,

		CNAME:       dns.Cname,
		ARecords:    dns.ARecords,
		AAAARecords: dns.AAAARecords,
		MXRecords:   dns.MxRecords,
		NSRecords:   dns.NsRecords,
		TXTRecords:  dns.TxtRecords,
		PTRRecord:   dns.PtrRecord,
		DNSSEC:      dns.Dnssec,
	}
}
