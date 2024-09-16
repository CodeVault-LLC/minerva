package rules

import (
	"github.com/codevault-llc/humblebrag-api/types"
	"github.com/codevault-llc/humblebrag-api/utils"
)

func TwitchAPIToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "twitch-api-token",
		Description: "Discovered a Twitch API token, which could compromise streaming services and account integrations.",
		Regex:       generateSemiGenericRegex([]string{"twitch"}, alphaNumeric("30"), true),
		Keywords: []string{
			"twitch",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("twitch", utils.NewSecret(alphaNumeric("30"))),
	}
	return validate(r, tps, nil)
}
