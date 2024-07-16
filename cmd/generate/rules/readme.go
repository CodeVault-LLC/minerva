package rules

import (
	"github.com/codevault-llc/humblebrag-api/cmd/generate/secrets"
	"github.com/codevault-llc/humblebrag-api/config"
)

func ReadMe() *config.Rule {
	// define rule
	r := config.Rule{
		RuleID:      "readme-api-token",
		Description: "Detected a Readme API token, risking unauthorized documentation management and content exposure.",
		Regex:       generateUniqueTokenRegex(`rdme_[a-z0-9]{70}`, true),

		Keywords: []string{
			"rdme_",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("api-token", "rdme_"+secrets.NewSecret(alphaNumeric("70"))),
	}
	return validate(r, tps, nil)
}
