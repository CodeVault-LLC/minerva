package network

import (
	"fmt"
	"net"

	"github.com/codevault-llc/humblebrag-api/utils"
)

// ScanIP scans the IP address of a given URL and returns a list of IP addresses
func ScanIP(url string) ([]string, error) {
	url = utils.StripProtocol(url)

	ips, err := net.LookupIP(url)
	if err != nil {
		fmt.Println("Error in IP Lookup:", err)
		return nil, err
	}

	var ipList []string
	for _, ip := range ips {
		ipList = append(ipList, ip.String())
	}

	return ipList, nil
}

// ScanIPRange scans the IP range of a given URL and returns a list of IP addresses
func ScanIPRange(url string) ([]string, error) {
	url = utils.StripProtocol(url)

	ips, err := net.LookupHost(url)
	if err != nil {
		return nil, err
	}

	return ips, nil
}
