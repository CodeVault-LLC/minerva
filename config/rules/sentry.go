package rules

import (
	"github.com/codevault-llc/humblebrag-api/types"
	"github.com/codevault-llc/humblebrag-api/utils"
)

func SentryAccessToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "sentry-access-token",
		Description: "Found a Sentry Access Token, risking unauthorized access to error tracking services and sensitive application data.",
		Regex:       generateSemiGenericRegex([]string{"sentry"}, hex("64"), true),

		Keywords: []string{
			"sentry",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("sentry", utils.NewSecret(hex("64"))),
	}
	return validate(r, tps, nil)
}
