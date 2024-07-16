package rules

import (
	"github.com/codevault-llc/humblebrag-api/cmd/generate/secrets"
	"github.com/codevault-llc/humblebrag-api/config"
)

func FinnhubAccessToken() *config.Rule {
	// define rule
	r := config.Rule{
		RuleID:      "finnhub-access-token",
		Description: "Found a Finnhub Access Token, risking unauthorized access to financial market data and analytics.",
		Regex:       generateSemiGenericRegex([]string{"finnhub"}, alphaNumeric("20"), true),

		Keywords: []string{
			"finnhub",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("finnhub", secrets.NewSecret(alphaNumeric("20"))),
	}
	return validate(r, tps, nil)
}
