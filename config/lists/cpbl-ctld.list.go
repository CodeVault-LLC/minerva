package lists

import (
	"strings"

	"github.com/codevault-llc/minerva/pkg/parsers"
)

var CblCtldParser = &parsers.TextParser{
	ParseFunc: func(line string) (parsers.Item, bool) {
		line = strings.TrimSpace(line)
		return parsers.Item{Type: parsers.Domain, Value: line}, true
	},
}
