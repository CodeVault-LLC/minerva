package lists

import (
	"strings"

	"github.com/codevault-llc/humblebrag-api/parsers"
)

/*
###
#
# Name: HOSTS Database
# Date: 06:01 AM Sunday 30 June 2024
#
# Repo: https://tgc.cloud/downloads
# Path: https://tgc.cloud/downloads/hosts.txt
#
#
###

127.0.0.1 localhost
255.255.255.255 broadcasthost
::1 localhost


0.0.0.0 spryptvqguhpwt
0.0.0.0 metrics.abbott
0.0.0.0 smetrics.abbott
0.0.0.0 marketing.globalpointofcare.abbott
0.0.0.0 marketing2.globalpointofcare.abbott
*/

var IPSumParser = &parsers.TextParser{
	ParseFunc: func(line string) (parsers.Item, bool) {
		line = strings.TrimSpace(line)

		// Ignore comment or metadata lines starting with '#'
		if strings.HasPrefix(line, "#") || line == "" {
			return parsers.Item{}, false
		}

		fields := strings.Fields(line)
		if len(fields) < 2 {
			return parsers.Item{}, false
		}

		ip := fields[0]
		domain := fields[1]

		// Check if line is a valid URL (starts with http/https)
		if isLocalhost(ip) || isLocalhost(domain) {
			return parsers.Item{}, false
		}

		return parsers.Item{Type: parsers.IPv4, Value: ip}, true
	},
}
