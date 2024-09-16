package rules

import (
	"github.com/codevault-llc/humblebrag-api/types"
	"github.com/codevault-llc/humblebrag-api/utils"
)

func GitterAccessToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "gitter-access-token",
		Description: "Uncovered a Gitter Access Token, which may lead to unauthorized access to chat and communication services.",
		Regex: generateSemiGenericRegex([]string{"gitter"},
			alphaNumericExtendedShort("40"), true),

		Keywords: []string{
			"gitter",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("gitter",
			utils.NewSecret(alphaNumericExtendedShort("40"))),
	}
	return validate(r, tps, nil)
}
