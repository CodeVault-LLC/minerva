package rules

import (
	"github.com/codevault-llc/minerva/pkg/types"
	"github.com/codevault-llc/minerva/pkg/utils"
)

func RapidAPIAccessToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
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
			utils.NewSecret(alphaNumericExtendedShort("50"))),
	}
	return validate(r, tps, nil)
}
