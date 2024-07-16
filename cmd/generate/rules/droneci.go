package rules

import (
	"github.com/codevault-llc/humblebrag-api/cmd/generate/secrets"
	"github.com/codevault-llc/humblebrag-api/config"
)

func DroneciAccessToken() *config.Rule {
	// define rule
	r := config.Rule{
		RuleID:      "droneci-access-token",
		Description: "Detected a Droneci Access Token, potentially compromising continuous integration and deployment workflows.",
		Regex:       generateSemiGenericRegex([]string{"droneci"}, alphaNumeric("32"), true),

		Keywords: []string{
			"droneci",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("droneci", secrets.NewSecret(alphaNumeric("32"))),
	}
	return validate(r, tps, nil)
}
