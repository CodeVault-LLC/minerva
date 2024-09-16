package rules

import (
	"github.com/codevault-llc/humblebrag-api/types"
	"github.com/codevault-llc/humblebrag-api/utils"
)

func OktaAccessToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
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
		generateSampleSecret("okta", utils.NewSecret(alphaNumeric("42"))),
	}
	return validate(r, tps, nil)
}
