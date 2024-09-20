package rules

import (
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
	"github.com/codevault-llc/humblebrag-api/types"
)

func TravisCIAccessToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "travisci-access-token",
		Description: "Identified a Travis CI Access Token, potentially compromising continuous integration services and codebase security.",
		Regex:       generateSemiGenericRegex([]string{"travis"}, alphaNumeric("22"), true),

		Keywords: []string{
			"travis",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("travis", utils.NewSecret(alphaNumeric("22"))),
	}
	return validate(r, tps, nil)
}
