package rules

import (
	regexp "github.com/wasilibs/go-re2"

	"github.com/codevault-llc/humblebrag-api/config"
)

func URLToken() *config.Rule {
	// define rule
	r := config.Rule{
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
