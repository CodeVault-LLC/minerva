package rules

import (
	"github.com/codevault-llc/minerva/pkg/types"
	"github.com/codevault-llc/minerva/pkg/utils"
)

func FinicityClientSecret() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Identified a Finicity Client Secret, which could lead to compromised financial service integrations and data breaches.",
		RuleID:      "finicity-client-secret",
		Regex:       generateSemiGenericRegex([]string{"finicity"}, alphaNumeric("20"), true),

		Keywords: []string{"finicity"},
	}

	// validate
	tps := []string{
		generateSampleSecret("finicity", utils.NewSecret(alphaNumeric("20"))),
	}
	return validate(r, tps, nil)
}

func FinicityAPIToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Detected a Finicity API token, potentially risking financial data access and unauthorized financial operations.",
		RuleID:      "finicity-api-token",
		Regex:       generateSemiGenericRegex([]string{"finicity"}, hex("32"), true),

		Keywords: []string{"finicity"},
	}

	// validate
	tps := []string{
		generateSampleSecret("finicity", utils.NewSecret(hex("32"))),
	}
	return validate(r, tps, nil)
}
