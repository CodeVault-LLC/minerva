package rules

import (
	"fmt"

	"github.com/codevault-llc/humblebrag-api/types"
	"github.com/codevault-llc/humblebrag-api/utils"
)

func PlaidAccessID() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "plaid-client-id",
		Description: "Uncovered a Plaid Client ID, which could lead to unauthorized financial service integrations and data breaches.",
		Regex:       generateSemiGenericRegex([]string{"plaid"}, alphaNumeric("24"), true),

		Keywords: []string{
			"plaid",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("plaid", utils.NewSecret(alphaNumeric("24"))),
	}
	return validate(r, tps, nil)
}

func PlaidSecretKey() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "plaid-secret-key",
		Description: "Detected a Plaid Secret key, risking unauthorized access to financial accounts and sensitive transaction data.",
		Regex:       generateSemiGenericRegex([]string{"plaid"}, alphaNumeric("30"), true),

		Keywords: []string{
			"plaid",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("plaid", utils.NewSecret(alphaNumeric("30"))),
	}
	return validate(r, tps, nil)
}

func PlaidAccessToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "plaid-api-token",
		Description: "Discovered a Plaid API Token, potentially compromising financial data aggregation and banking services.",
		Regex: generateSemiGenericRegex([]string{"plaid"},
			fmt.Sprintf("access-(?:sandbox|development|production)-%s", hex8_4_4_4_12()), true),

		Keywords: []string{
			"plaid",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("plaid", utils.NewSecret(fmt.Sprintf("access-(?:sandbox|development|production)-%s", hex8_4_4_4_12()))),
	}
	return validate(r, tps, nil)
}
