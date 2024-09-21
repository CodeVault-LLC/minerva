package rules

import (
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
	"github.com/codevault-llc/humblebrag-api/types"
)

func Prefect() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "prefect-api-token",
		Description: "Detected a Prefect API token, risking unauthorized access to workflow management and automation services.",
		Regex:       generateUniqueTokenRegex(`pnu_[a-z0-9]{36}`, true),

		Keywords: []string{
			"pnu_",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("api-token", "pnu_"+utils.NewSecret(alphaNumeric("36"))),
	}
	return validate(r, tps, nil)
}
