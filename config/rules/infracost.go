package rules

import (
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
	"github.com/codevault-llc/humblebrag-api/types"
)

func InfracostAPIToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
		// Human readable description of the types.Rule
		Description: "Detected an Infracost API Token, risking unauthorized access to cloud cost estimation tools and financial data.",

		// Unique ID for the types.Rule
		RuleID: "infracost-api-token",

		// Regex capture group for the actual secret

		// Regex used for detecting secrets. See regex section below for more details
		Regex: generateUniqueTokenRegex(`ico-[a-zA-Z0-9]{32}`, true),

		// Keywords used for string matching on fragments (think of this as a prefilter)
		Keywords: []string{"ico-"},
	}

	// validate
	tps := []string{
		generateSampleSecret("ico", "ico-"+utils.NewSecret("[A-Za-z0-9]{32}")),
	}
	return validate(r, tps, nil)
}
