package types

import "github.com/codevault-llc/humblebrag-api/pkg/parsers"

type Filter struct {
	Description string
	FilterID      string
	Categories  []string
	Types       []parsers.ListType
	URL         string
	Parser      parsers.Parser
}
