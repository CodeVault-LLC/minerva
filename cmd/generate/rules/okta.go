package rules

import (
	"github.com/codevault-llc/humblebrag-api/cmd/generate/secrets"
	"github.com/codevault-llc/humblebrag-api/config"
)

func OktaAccessToken() *config.Rule {
	// define rule
	r := config.Rule{
		RuleID:      "okta-access-token",
		Description: "Identified an Okta Access Token, which may compromise identity management services and user authentication data.",
		Regex: generateSemiGenericRegex([]string{"okta"},
			alphaNumericExtended("42"), true),

		Keywords: []string{
			"okta",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("okta", secrets.NewSecret(alphaNumeric("42"))),
	}
	return validate(r, tps, nil)
}
