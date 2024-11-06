package rules

import (
	"github.com/codevault-llc/minerva/pkg/types"
	"github.com/codevault-llc/minerva/pkg/utils"
)

func CoinbaseAccessToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "coinbase-access-token",
		Description: "Detected a Coinbase Access Token, posing a risk of unauthorized access to cryptocurrency accounts and financial transactions.",
		Regex: generateSemiGenericRegex([]string{"coinbase"},
			alphaNumericExtendedShort("64"), true),
		Keywords: []string{
			"coinbase",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("coinbase",
			utils.NewSecret(alphaNumericExtendedShort("64"))),
	}
	return validate(r, tps, nil)
}
