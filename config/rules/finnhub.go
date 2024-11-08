package rules

import (
	"github.com/codevault-llc/minerva/pkg/types"
	"github.com/codevault-llc/minerva/pkg/utils"
)

func FinnhubAccessToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "finnhub-access-token",
		Description: "Found a Finnhub Access Token, risking unauthorized access to financial market data and analytics.",
		Regex:       generateSemiGenericRegex([]string{"finnhub"}, alphaNumeric("20"), true),

		Keywords: []string{
			"finnhub",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("finnhub", utils.NewSecret(alphaNumeric("20"))),
	}
	return validate(r, tps, nil)
}
