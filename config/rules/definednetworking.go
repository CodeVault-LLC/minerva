package rules

import (
	"github.com/codevault-llc/humblebrag-api/pkg/types"
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
)

func DefinedNetworkingAPIToken() *types.Rule {
	// Define types.Rule
	r := types.Rule{
		// Human redable description of the types.Rule
		Description: "Identified a Defined Networking API token, which could lead to unauthorized network operations and data breaches.",

		// Unique ID for the types.Rule
		RuleID: "defined-networking-api-token",

		// Regex used for detecting secrets. See regex section below for more details
		Regex: generateSemiGenericRegex([]string{"dnkey"}, `dnkey-[a-z0-9=_\-]{26}-[a-z0-9=_\-]{52}`, true),

		// Keywords used for string matching on fragments (think of this as a prefilter)
		Keywords: []string{"dnkey"},
	}

	// validate
	tps := []string{
		generateSampleSecret("dnkey", "dnkey-"+utils.NewSecret(alphaNumericExtended("26"))+"-"+utils.NewSecret(alphaNumericExtended("52"))),
	}
	return validate(r, tps, nil)
}
