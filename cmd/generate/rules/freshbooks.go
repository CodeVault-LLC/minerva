package rules

import (
	"github.com/codevault-llc/humblebrag-api/cmd/generate/secrets"
	"github.com/codevault-llc/humblebrag-api/config"
)

func FreshbooksAccessToken() *config.Rule {
	// define rule
	r := config.Rule{
		RuleID:      "freshbooks-access-token",
		Description: "Discovered a Freshbooks Access Token, posing a risk to accounting software access and sensitive financial data exposure.",
		Regex:       generateSemiGenericRegex([]string{"freshbooks"}, alphaNumeric("64"), true),

		Keywords: []string{
			"freshbooks",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("freshbooks", secrets.NewSecret(alphaNumeric("64"))),
	}
	return validate(r, tps, nil)
}
