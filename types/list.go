package types

import "github.com/codevault-llc/humblebrag-api/parsers"

type List struct {
	Description string
	ListID      string
	Categories  []string
	URL         string
	Parser      parsers.Parser
}
