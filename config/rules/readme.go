package rules

import (
	"github.com/codevault-llc/minerva/pkg/types"
	"github.com/codevault-llc/minerva/pkg/utils"
)

func ReadMe() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "readme-api-token",
		Description: "Detected a Readme API token, risking unauthorized documentation management and content exposure.",
		Regex:       generateUniqueTokenRegex(`rdme_[a-z0-9]{70}`, true),

		Keywords: []string{
			"rdme_",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("api-token", "rdme_"+utils.NewSecret(alphaNumeric("70"))),
	}
	return validate(r, tps, nil)
}
