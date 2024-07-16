package rules

import (
	"github.com/codevault-llc/humblebrag-api/cmd/generate/secrets"
	"github.com/codevault-llc/humblebrag-api/config"
)

func MattermostAccessToken() *config.Rule {
	// define rule
	r := config.Rule{
		RuleID:      "mattermost-access-token",
		Description: "Identified a Mattermost Access Token, which may compromise team communication channels and data privacy.",
		Regex:       generateSemiGenericRegex([]string{"mattermost"}, alphaNumeric("26"), true),

		Keywords: []string{
			"mattermost",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("mattermost", secrets.NewSecret(alphaNumeric("26"))),
	}
	return validate(r, tps, nil)
}
