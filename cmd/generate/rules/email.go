package rules

import (
	regexp "github.com/wasilibs/go-re2"

	"github.com/codevault-llc/humblebrag-api/config"
)

func EmailToken() *config.Rule {
	// define rule
	r := config.Rule{
		Description: "Identifies Emails in the codebase.",
		RuleID:      "email-token",
		Regex:       regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`),
		Keywords:    []string{"email", "gmail", "@"},
	}

	// validate
	tps := []string{
		"fakeemail@gmail.com",
		"thisisadomainemail@domain.nq",
	}

	return validate(r, tps, nil)
}
