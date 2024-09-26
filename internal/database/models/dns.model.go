package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type DNSModel struct {
	gorm.Model

	NetworkId uint
	Network   *NetworkModel `gorm:"foreignKey:NetworkId"`

	CNAME       pq.StringArray `gorm:"type:text[]"` // PostgreSQL array
	ARecords    pq.StringArray `gorm:"type:text[]"` // PostgreSQL array
	AAAARecords pq.StringArray `gorm:"type:text[]"` // PostgreSQL array
	MXRecords   pq.StringArray `gorm:"type:text[]"` // PostgreSQL array
	NSRecords   pq.StringArray `gorm:"type:text[]"` // PostgreSQL array
	TXTRecords  pq.StringArray `gorm:"type:text[]"` // PostgreSQL array
	PTRRecord   string         `gorm:"not null"`
	DNSSEC      bool           `gorm:"not null"`
}

type DNSResponse struct {
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

func ConvertDNS(dns DNSModel) DNSResponse {
	return DNSResponse{
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
