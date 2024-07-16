package rules

import (
	"github.com/codevault-llc/humblebrag-api/cmd/generate/secrets"
	"github.com/codevault-llc/humblebrag-api/config"
)

func NetlifyAccessToken() *config.Rule {
	// define rule
	r := config.Rule{
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
		generateSampleSecret("netlify", secrets.NewSecret(alphaNumericExtended("40,46"))),
	}
	return validate(r, tps, nil)
}
