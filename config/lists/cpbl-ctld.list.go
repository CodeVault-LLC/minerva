package lists

import (
	"strings"

	"github.com/codevault-llc/humblebrag-api/parsers"
)

var URLHausParser = &parsers.TextParser{
	ParseFunc: func(line string) (parsers.Item, bool) {
		line = strings.TrimSpace(line)

		// Ignore comment or metadata lines starting with '#'
		if strings.HasPrefix(line, "#") || line == "" {
			return parsers.Item{}, false
		}

		// Check if line is a valid URL (starts with http/https)
		if strings.HasPrefix(line, "http://") || strings.HasPrefix(line, "https://") {
			return parsers.Item{Type: "URL", Value: line}, true
		}

		return parsers.Item{}, false
	},
}
