package rules

import (
	"github.com/codevault-llc/humblebrag-api/cmd/generate/secrets"
	"github.com/codevault-llc/humblebrag-api/config"
)

func TwitchAPIToken() *config.Rule {
	// define rule
	r := config.Rule{
		RuleID:      "twitch-api-token",
		Description: "Discovered a Twitch API token, which could compromise streaming services and account integrations.",
		Regex:       generateSemiGenericRegex([]string{"twitch"}, alphaNumeric("30"), true),
		Keywords: []string{
			"twitch",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("twitch", secrets.NewSecret(alphaNumeric("30"))),
	}
	return validate(r, tps, nil)
}
