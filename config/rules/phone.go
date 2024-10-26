package rules

import (
	"github.com/codevault-llc/humblebrag-api/pkg/types"
	regexp "github.com/wasilibs/go-re2"
)

func PhoneToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Identifies phone numbers in the codebase.",
		RuleID:      "phone-token",
		Regex:       regexp.MustCompile(`\+?[0-9]{1,3}[-\s]?\(?[0-9]{3}\)?[-\s]?[0-9]{3}[-\s]?[0-9]{4}`),
		Keywords:    []string{"phone", "number"},
	}

	// validate
	tps := []string{
		"+1-555-555-5555",
		"555-555-5555",
		"5555555555",
		"+1 555 555 5555",
	}

	// non matches
	fps := []string{
		"555-555-555",
		"555-555-55555",
	}

	return validate(r, tps, fps)
}
