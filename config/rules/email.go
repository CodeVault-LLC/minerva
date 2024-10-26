package rules

import (
	"github.com/codevault-llc/humblebrag-api/pkg/types"
	regexp "github.com/wasilibs/go-re2"
)

func EmailToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
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
