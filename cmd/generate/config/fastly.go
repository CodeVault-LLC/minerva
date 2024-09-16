package rules

import (
	"github.com/codevault-llc/humblebrag-api/cmd/generate/secrets"
	"github.com/codevault-llc/humblebrag-api/config"
)

func FastlyAPIToken() *config.Rule {
	// define rule
	r := config.Rule{
		Description: "Uncovered a Fastly API key, which may compromise CDN and edge cloud services, leading to content delivery and security issues.",
		RuleID:      "fastly-api-token",
		Regex:       generateSemiGenericRegex([]string{"fastly"}, alphaNumericExtended("32"), true),

		Keywords: []string{"fastly"},
	}

	// validate
	tps := []string{
		generateSampleSecret("fastly", secrets.NewSecret(alphaNumericExtended("32"))),
	}
	return validate(r, tps, nil)
}
