package rules

import (
	"github.com/codevault-llc/minerva/pkg/types"
	"github.com/codevault-llc/minerva/pkg/utils"
)

func DatadogtokenAccessToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "datadog-access-token",
		Description: "Detected a Datadog Access Token, potentially risking monitoring and analytics data exposure and manipulation.",
		Regex: generateSemiGenericRegex([]string{"datadog"},
			alphaNumeric("40"), true),
		Keywords: []string{
			"datadog",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("datadog", utils.NewSecret(alphaNumeric("40"))),
	}
	return validate(r, tps, nil)
}
