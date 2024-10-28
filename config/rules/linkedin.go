package rules

import (
	"github.com/codevault-llc/humblebrag-api/pkg/types"
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
)

func LinkedinClientSecret() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "linkedin-client-secret",
		Description: "Discovered a LinkedIn Client secret, potentially compromising LinkedIn application integrations and user data.",
		Regex: generateSemiGenericRegex([]string{
			"linkedin",
			"linked-in",
		}, alphaNumeric("16"), true),

		Keywords: []string{
			"linkedin",
			"linked-in",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("linkedin", utils.NewSecret(alphaNumeric("16"))),
	}
	return validate(r, tps, nil)
}

func LinkedinClientID() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "linkedin-client-id",
		Description: "Found a LinkedIn Client ID, risking unauthorized access to LinkedIn integrations and professional data exposure.",
		Regex: generateSemiGenericRegex([]string{
			"linkedin",
			"linked-in",
		}, alphaNumeric("14"), true),

		Keywords: []string{
			"linkedin",
			"linked-in",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("linkedin", utils.NewSecret(alphaNumeric("14"))),
	}
	return validate(r, tps, nil)
}
