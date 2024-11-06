package entities

import (
	"crypto/x509"
	"time"

	"github.com/lib/pq"
)

type NetworkModel struct {
	Id uint `db:"id"`

	ScanId uint `db:"scan_id"`

	IpAddresses pq.StringArray `db:"ip_addresses"`
	IpRanges    pq.StringArray `db:"ip_ranges"`

	HttpHeaders pq.StringArray `db:"http_headers"`

	// Relations
	Certificates []CertificateModel `db:"-"`
	Whois        WhoisModel         `db:"-"`
	DNS          DnsModel           `db:"-"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type WhoisModel struct {
	Id uint `db:"id"`

	NetworkId uint `db:"network_id"`

	DomainName  string         `db:"domain_name"`
	Registrar   string         `db:"registrar"`
	Email       string         `db:"email"`
	Phone       string         `db:"phone"`
	Updated     string         `db:"updated"`
	Created     string         `db:"created"`
	Expires     string         `db:"expires"`
	Status      string         `db:"status"`
	NameServers pq.StringArray `db:"name_servers"`

	RegistrantName       string `db:"registrant_name"`
	RegistrantEmail      string `db:"registrant_email"`
	RegistrantPhone      string `db:"registrant_phone"`
	RegistrantOrg        string `db:"registrant_org"`
	RegistrantCity       string `db:"registrant_city"`
	RegistrantCountry    string `db:"registrant_country"`
	RegistrantPostalCode string `db:"registrant_postal_code"`

	AdminName       string `db:"admin_name"`
	AdminEmail      string `db:"admin_email"`
	AdminPhone      string `db:"admin_phone"`
	AdminOrg        string `db:"admin_org"`
	AdminCity       string `db:"admin_city"`
	AdminCountry    string `db:"admin_country"`
	AdminPostalCode string `db:"admin_postal_code"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type DnsModel struct {
	Id uint `db:"id"`

	NetworkId uint `db:"network_id"`

	Cname       pq.StringArray `db:"cname"`
	ARecords    pq.StringArray `db:"a_records"`
	AAAARecords pq.StringArray `db:"aaaa_records"`
	MxRecords   pq.StringArray `db:"mx_records"`
	NsRecords   pq.StringArray `db:"ns_records"`
	TxtRecords  pq.StringArray `db:"txt_records"`
	PtrRecord   string         `db:"ptr_record"`
	Dnssec      bool           `db:"dnssec"`
}

type CertificateModel struct {
	Id uint `db:"id"`

	NetworkId uint `db:"network_id"`

	Subject string `db:"subject"`
	Issuer  string `db:"issuer"`

	NotBefore time.Time `db:"not_before"`
	NotAfter  time.Time `db:"not_after"`

	SignatureAlgorithm x509.SignatureAlgorithm `db:"signature_algorithm"`
	Signature          []byte                  `db:"signature"`

	PublicKeyAlgorithm x509.PublicKeyAlgorithm `db:"public_key_algorithm"`
	PublicKey          string                  `db:"public_key"`

	SerialNumber string        `db:"serial_number"`
	Version      int           `db:"version"`
	KeyUsage     x509.KeyUsage `db:"key_usage"`

	BasicConstraintsValid bool `db:"basic_constraints_valid"`
	IsCa                  bool `db:"is_ca"`

	DnsNames       pq.StringArray `db:"dns_names"`
	EmailAddresses pq.StringArray `db:"email_addresses"`
	IpAddresses    pq.StringArray `db:"ip_addresses"`
	Uris           pq.StringArray `db:"uris"`

	PermittedDnsDomainsCritical bool     `db:"permitted_dns_domains_critical"`
	PermittedDnsDomains         pq.StringArray `db:"permitted_dns_domains"`
	ExcludedDnsDomains          pq.StringArray `db:"excluded_dns_domains"`
	PermittedIpRanges           pq.StringArray `db:"permitted_ip_ranges"`
	ExcludedIpRanges            pq.StringArray `db:"excluded_ip_ranges"`
	PermittedEmailAddresses     pq.StringArray `db:"permitted_email_addresses"`
	ExcludedEmailAddresses      pq.StringArray `db:"excluded_email_addresses"`
	PermittedUriDomains         pq.StringArray `db:"permitted_uri_domains"`
	ExcludedUriDomains          pq.StringArray `db:"excluded_uri_domains"`
}
