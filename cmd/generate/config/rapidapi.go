package rules

import (
	"github.com/codevault-llc/humblebrag-api/cmd/generate/secrets"
	"github.com/codevault-llc/humblebrag-api/config"
)

func RapidAPIAccessToken() *config.Rule {
	// define rule
	r := config.Rule{
		RuleID:      "rapidapi-access-token",
		Description: "Uncovered a RapidAPI Access Token, which could lead to unauthorized access to various APIs and data services.",
		Regex: generateSemiGenericRegex([]string{"rapidapi"},
			alphaNumericExtendedShort("50"), true),

		Keywords: []string{
			"rapidapi",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("rapidapi",
			secrets.NewSecret(alphaNumericExtendedShort("50"))),
	}
	return validate(r, tps, nil)
}
