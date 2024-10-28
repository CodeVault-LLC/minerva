package rules

import (
	"github.com/codevault-llc/humblebrag-api/pkg/types"
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
)

func SquareSpaceAccessToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "squarespace-access-token",
		Description: "Identified a Squarespace Access Token, which may compromise website management and content control on Squarespace.",
		Regex:       generateSemiGenericRegex([]string{"squarespace"}, hex8_4_4_4_12(), true),

		Keywords: []string{
			"squarespace",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("squarespace", utils.NewSecret(hex8_4_4_4_12())),
	}
	return validate(r, tps, nil)
}
