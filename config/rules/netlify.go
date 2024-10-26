package rules

import (
	"github.com/codevault-llc/humblebrag-api/pkg/types"
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
)

func NetlifyAccessToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "netlify-access-token",
		Description: "Detected a Netlify Access Token, potentially compromising web hosting services and site management.",
		Regex: generateSemiGenericRegex([]string{"netlify"},
			alphaNumericExtended("40,46"), true),

		Keywords: []string{
			"netlify",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("netlify", utils.NewSecret(alphaNumericExtended("40,46"))),
	}
	return validate(r, tps, nil)
}
