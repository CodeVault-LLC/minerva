package lists

import (
	"strings"

	"github.com/codevault-llc/humblebrag-api/parsers"
)

var CblAbpParser = &parsers.TextParser{
	ParseFunc: func(line string) (parsers.Item, bool) {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "||") && strings.HasSuffix(line, "^") {
			domain := strings.TrimPrefix(strings.TrimSuffix(line, "^"), "||")
			return parsers.Item{Type: "Domain", Value: domain}, true
		}
		return parsers.Item{}, false
	},
}
