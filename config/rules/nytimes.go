package rules

import (
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
	"github.com/codevault-llc/humblebrag-api/types"
)

func NytimesAccessToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "nytimes-access-token",
		Description: "Detected a Nytimes Access Token, risking unauthorized access to New York Times APIs and content services.",
		Regex: generateSemiGenericRegex([]string{
			"nytimes", "new-york-times,", "newyorktimes"},
			alphaNumericExtended("32"), true),

		Keywords: []string{
			"nytimes",
			"new-york-times",
			"newyorktimes",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("nytimes", utils.NewSecret(alphaNumeric("32"))),
	}
	return validate(r, tps, nil)
}
