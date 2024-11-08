package rules

import (
	"github.com/codevault-llc/minerva/pkg/types"
	"github.com/codevault-llc/minerva/pkg/utils"
)

func FastlyAPIToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Uncovered a Fastly API key, which may compromise CDN and edge cloud services, leading to content delivery and security issues.",
		RuleID:      "fastly-api-token",
		Regex:       generateSemiGenericRegex([]string{"fastly"}, alphaNumericExtended("32"), true),

		Keywords: []string{"fastly"},
	}

	// validate
	tps := []string{
		generateSampleSecret("fastly", utils.NewSecret(alphaNumericExtended("32"))),
	}
	return validate(r, tps, nil)
}
