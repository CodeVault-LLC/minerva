package rules

import (
	"github.com/codevault-llc/humblebrag-api/cmd/generate/secrets"
	"github.com/codevault-llc/humblebrag-api/config"
)

func NPM() *config.Rule {
	// define rule
	r := config.Rule{
		RuleID:      "npm-access-token",
		Description: "Uncovered an npm access token, potentially compromising package management and code repository access.",
		Regex:       generateUniqueTokenRegex(`npm_[a-z0-9]{36}`, true),

		Keywords: []string{
			"npm_",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("npmAccessToken", "npm_"+secrets.NewSecret(alphaNumeric("36"))),
	}
	return validate(r, tps, nil)
}
