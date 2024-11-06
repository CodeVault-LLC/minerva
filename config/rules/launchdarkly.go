package rules

import (
	"github.com/codevault-llc/minerva/pkg/types"
	"github.com/codevault-llc/minerva/pkg/utils"
)

func LaunchDarklyAccessToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "launchdarkly-access-token",
		Description: "Uncovered a Launchdarkly Access Token, potentially compromising feature flag management and application functionality.",
		Regex:       generateSemiGenericRegex([]string{"launchdarkly"}, alphaNumericExtended("40"), true),

		Keywords: []string{
			"launchdarkly",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("launchdarkly", utils.NewSecret(alphaNumericExtended("40"))),
	}
	return validate(r, tps, nil)
}
