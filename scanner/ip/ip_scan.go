package ip

import "net"

func ScanIP(url string) ([]string, error) {
	ips, err := net.LookupIP(url)
	if err != nil {
		return nil, err
	}

	var ipList []string
	for _, ip := range ips {
		ipList = append(ipList, ip.String())
	}

	return ipList, nil
}

func ScanIPRange(url string) ([]string, error) {
	ips, err := net.LookupHost(url)
	if err != nil {
		return nil, err
	}

	return ips, nil
}
