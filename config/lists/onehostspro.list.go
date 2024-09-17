package lists

import (
	"strings"

	"github.com/codevault-llc/humblebrag-api/parsers"
)

var OneHostsProParser = &parsers.TextParser{
	ParseFunc: func(line string) (parsers.Item, bool) {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "#") || strings.HasPrefix(line, "/") || strings.HasPrefix(line, "\\") || strings.HasPrefix(line, "(") || strings.HasPrefix(line, "|") {
			return parsers.Item{}, false
		}

		fields := strings.Fields(line)
		if len(fields) < 2 {
			return parsers.Item{}, false
		}

		ip := fields[0]
		domain := fields[1]

		if isLocalhost(ip) || isLocalhost(domain) {
			return parsers.Item{}, false
		}

		return parsers.Item{Type: parsers.Domain, Value: domain}, true
	},
}
