package entities

import (
	"crypto/x509"
	"time"

	whoisparser "github.com/likexian/whois-parser"
)

type NetworkModel struct {
	Id     string `cql:"id"`
	ScanId string `cql:"scan_id"`

	IpAddresses []string `cql:"ip_addresses"`
	IpRanges    []string `cql:"ip_ranges"`

	HttpHeaders      []string              `cql:"http_headers"`
	WhoisModel       whoisparser.WhoisInfo `cql:"whois"`
	DnsModel         DnsModel              `cql:"dns"`
	CertificateModel []*x509.Certificate   `cql:"certificates"`

	CreatedAt time.Time `cql:"created_at"`
	UpdatedAt time.Time `cql:"updated_at"`
}

type WhoisModel struct {
	DomainName  string   `cql:"domain_name"`
	Registrar   string   `cql:"registrar"`
	Email       string   `cql:"email"`
	Phone       string   `cql:"phone"`
	Updated     string   `cql:"updated"`
	Created     string   `cql:"created"`
	Expires     string   `cql:"expires"`
	Status      string   `cql:"status"`
	NameServers []string `cql:"name_servers"`

	RegistrantName       string `cql:"registrant_name"`
	RegistrantEmail      string `cql:"registrant_email"`
	RegistrantPhone      string `cql:"registrant_phone"`
	RegistrantOrg        string `cql:"registrant_org"`
	RegistrantCity       string `cql:"registrant_city"`
	RegistrantCountry    string `cql:"registrant_country"`
	RegistrantPostalCode string `cql:"registrant_postal_code"`

	AdminName       string `cql:"admin_name"`
	AdminEmail      string `cql:"admin_email"`
	AdminPhone      string `cql:"admin_phone"`
	AdminOrg        string `cql:"admin_org"`
	AdminCity       string `cql:"admin_city"`
	AdminCountry    string `cql:"admin_country"`
	AdminPostalCode string `cql:"admin_postal_code"`

	CreatedAt time.Time `cql:"created_at"`
	UpdatedAt time.Time `cql:"updated_at"`
}

type DnsModel struct {
	Cname       []string `cql:"cname"`
	ARecords    []string `cql:"a_records"`
	AAAARecords []string `cql:"aaaa_records"`
	MxRecords   []string `cql:"mx_records"`
	NsRecords   []string `cql:"ns_records"`
	TxtRecords  []string `cql:"txt_records"`
	PtrRecord   string   `cql:"ptr_record"`
	Dnssec      bool     `cql:"dnssec"`
}
