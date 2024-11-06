package rules

import (
	"github.com/codevault-llc/minerva/pkg/types"
	"github.com/codevault-llc/minerva/pkg/utils"
)

func MattermostAccessToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "mattermost-access-token",
		Description: "Identified a Mattermost Access Token, which may compromise team communication channels and data privacy.",
		Regex:       generateSemiGenericRegex([]string{"mattermost"}, alphaNumeric("26"), true),

		Keywords: []string{
			"mattermost",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("mattermost", utils.NewSecret(alphaNumeric("26"))),
	}
	return validate(r, tps, nil)
}
