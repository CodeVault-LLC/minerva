package rules

import (
	"github.com/codevault-llc/minerva/pkg/types"
	"github.com/codevault-llc/minerva/pkg/utils"
)

func Databricks() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Uncovered a Databricks API token, which may compromise big data analytics platforms and sensitive data processing.",
		RuleID:      "databricks-api-token",
		Regex:       generateUniqueTokenRegex(`dapi[a-h0-9]{32}`, true),
		Keywords:    []string{"dapi"},
	}

	// validate
	tps := []string{
		generateSampleSecret("databricks", "dapi"+utils.NewSecret(hex("32"))),
	}
	return validate(r, tps, nil)
}
