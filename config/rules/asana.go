package rules

import (
	"github.com/codevault-llc/humblebrag-api/types"
	"github.com/codevault-llc/humblebrag-api/utils"
)

func AsanaClientID() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Discovered a potential Asana Client ID, risking unauthorized access to Asana projects and sensitive task information.",
		RuleID:      "asana-client-id",
		Regex:       generateSemiGenericRegex([]string{"asana"}, numeric("16"), true),
		Keywords:    []string{"asana"},
	}

	// validate
	tps := []string{
		generateSampleSecret("asana", utils.NewSecret(numeric("16"))),
	}
	return validate(r, tps, nil)
}

func AsanaClientSecret() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Identified an Asana Client Secret, which could lead to compromised project management integrity and unauthorized access.",
		RuleID:      "asana-client-secret",
		Regex:       generateSemiGenericRegex([]string{"asana"}, alphaNumeric("32"), true),

		Keywords: []string{"asana"},
	}

	// validate
	tps := []string{
		generateSampleSecret("asana", utils.NewSecret(alphaNumeric("32"))),
	}
	return validate(r, tps, nil)
}
