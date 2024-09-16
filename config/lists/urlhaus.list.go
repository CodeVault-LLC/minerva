package lists

import (
	"strings"

	"github.com/codevault-llc/humblebrag-api/parsers"
)

var CblCtldParser = &parsers.TextParser{
	ParseFunc: func(line string) (parsers.Item, bool) {
		line = strings.TrimSpace(line)
		return parsers.Item{Type: "Domain", Value: line}, true
	},
}
