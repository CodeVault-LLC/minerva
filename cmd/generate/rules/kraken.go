package rules

import (
	"github.com/codevault-llc/humblebrag-api/cmd/generate/secrets"
	"github.com/codevault-llc/humblebrag-api/config"
)

func KrakenAccessToken() *config.Rule {
	// define rule
	r := config.Rule{
		RuleID:      "kraken-access-token",
		Description: "Identified a Kraken Access Token, potentially compromising cryptocurrency trading accounts and financial security.",
		Regex: generateSemiGenericRegex([]string{"kraken"},
			alphaNumericExtendedLong("80,90"), true),

		Keywords: []string{
			"kraken",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("kraken",
			secrets.NewSecret(alphaNumericExtendedLong("80,90"))),
	}
	return validate(r, tps, nil)
}
