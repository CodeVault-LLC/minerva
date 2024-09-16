package config

import "github.com/codevault-llc/humblebrag-api/parsers"

type List struct {
	Description string
	ListID      string
	Categories  []string
	URL         string
	ParserID    string
}

func FindParser(parserID string) parsers.Parser {
	switch parserID {
	case "text":
		return &parsers.TextParser{}
	default:
		return nil
	}
}
