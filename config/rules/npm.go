package rules

import (
	"github.com/codevault-llc/minerva/pkg/types"
	"github.com/codevault-llc/minerva/pkg/utils"
)

func NPM() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "npm-access-token",
		Description: "Uncovered an npm access token, potentially compromising package management and code repository access.",
		Regex:       generateUniqueTokenRegex(`npm_[a-z0-9]{36}`, true),

		Keywords: []string{
			"npm_",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("npmAccessToken", "npm_"+utils.NewSecret(alphaNumeric("36"))),
	}
	return validate(r, tps, nil)
}
