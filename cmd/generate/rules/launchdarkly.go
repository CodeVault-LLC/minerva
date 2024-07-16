package rules

import (
	"github.com/codevault-llc/humblebrag-api/cmd/generate/secrets"
	"github.com/codevault-llc/humblebrag-api/config"
)

func LaunchDarklyAccessToken() *config.Rule {
	// define rule
	r := config.Rule{
		RuleID:      "launchdarkly-access-token",
		Description: "Uncovered a Launchdarkly Access Token, potentially compromising feature flag management and application functionality.",
		Regex:       generateSemiGenericRegex([]string{"launchdarkly"}, alphaNumericExtended("40"), true),

		Keywords: []string{
			"launchdarkly",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("launchdarkly", secrets.NewSecret(alphaNumericExtended("40"))),
	}
	return validate(r, tps, nil)
}
