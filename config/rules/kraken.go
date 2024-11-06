package rules

import (
	"github.com/codevault-llc/minerva/pkg/types"
	"github.com/codevault-llc/minerva/pkg/utils"
)

func KrakenAccessToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
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
			utils.NewSecret(alphaNumericExtendedLong("80,90"))),
	}
	return validate(r, tps, nil)
}
