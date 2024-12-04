package viewmodels

import "github.com/codevault-llc/minerva/internal/network/models/entities"

type Network struct {
	ID uint `json:"id"`

	IPAddresses []string `json:"ip_addresses"`
	IPRanges    []string `json:"ip_ranges"`

	HTTPHeaders []string `json:"http_headers"`

	Whois Whois `json:"whois"`
	DNS   DNS   `json:"dns"`
}

func ConvertNetwork(network entities.NetworkModel) Network {
	return Network{
		ID:          network.Id,
		IPAddresses: network.IpAddresses,
		IPRanges:    network.IpRanges,
		HTTPHeaders: network.HttpHeaders,
	}
}
