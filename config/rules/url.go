package rules

import (
	"github.com/codevault-llc/humblebrag-api/pkg/types"
	regexp "github.com/wasilibs/go-re2"
)

func URLToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Identifies different kinds of URLs in the codebase.",
		RuleID:      "url-token",
		Regex:       regexp.MustCompile(`(http|https):\/\/[a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,}(\/\S*)?`),
		Keywords:    []string{"http", "https"},
	}

	// validate
	tps := []string{
		"http://www.google.com",
		"https://www.google.com",
	}

	return validate(r, tps, nil)
}
