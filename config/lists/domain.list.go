package lists

import (
	"strings"

	"github.com/codevault-llc/minerva/pkg/parsers"
)

/*
# Title: HaGeZi's Ultimate DNS Blocklist
# Description: Ultimate Sweeper - Strictly cleans the Internet and protects your privacy! Blocks Ads, Affiliate, Tracking, Metrics, Telemetry, Phishing, Malware, Scam, Free Hoster, Fake, Coins and other "Crap".
# Homepage: https://github.com/hagezi/dns-blocklists
# License: https://github.com/hagezi/dns-blocklists/blob/main/LICENSE
# Issues: https://github.com/hagezi/dns-blocklists/issues
# Expires: 1 day
# Last modified: 17 Sep 2024 14:37 UTC
# Version: 2024.0917.1437.19
# Syntax: Domains (including possible subdomains)
# Number of entries: 668220
#
0.org
0.to
ellas2.0.org
www.0.org
www.0.to
0--4.com
www.0--4.com
0--d.com
*/

var DomainListParser = &parsers.TextParser{
	ParseFunc: func(line string) (parsers.Item, bool) {
		line = strings.TrimSpace(line)

		// Ignore comment or metadata lines starting with '#'
		if strings.HasPrefix(line, "#") || line == "" {
			return parsers.Item{}, false
		}

		// Remove www. from the URL to ensure consistency
		line = strings.Replace(line, "www.", "", 1)

		return parsers.Item{Type: parsers.Domain, Value: line}, true
	},
}
