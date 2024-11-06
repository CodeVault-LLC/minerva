package rules

import (
	"github.com/codevault-llc/minerva/pkg/types"
	"github.com/codevault-llc/minerva/pkg/utils"
)

func Typeform() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "typeform-api-token",
		Description: "Uncovered a Typeform API token, which could lead to unauthorized survey management and data collection.",
		Regex: generateSemiGenericRegex([]string{"typeform"},
			`tfp_[a-z0-9\-_\.=]{59}`, true),
		Keywords: []string{
			"tfp_",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("typeformAPIToken", "tfp_"+utils.NewSecret(alphaNumericExtended("59"))),
	}
	return validate(r, tps, nil)
}
