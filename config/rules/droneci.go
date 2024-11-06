package rules

import (
	"github.com/codevault-llc/minerva/pkg/types"
	"github.com/codevault-llc/minerva/pkg/utils"
)

func DroneciAccessToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "droneci-access-token",
		Description: "Detected a Droneci Access Token, potentially compromising continuous integration and deployment workflows.",
		Regex:       generateSemiGenericRegex([]string{"droneci"}, alphaNumeric("32"), true),

		Keywords: []string{
			"droneci",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("droneci", utils.NewSecret(alphaNumeric("32"))),
	}
	return validate(r, tps, nil)
}
