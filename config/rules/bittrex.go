package rules

import (
	"github.com/codevault-llc/humblebrag-api/pkg/types"
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
)

func BittrexAccessKey() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Identified a Bittrex Access Key, which could lead to unauthorized access to cryptocurrency trading accounts and financial loss.",
		RuleID:      "bittrex-access-key",
		Regex:       generateSemiGenericRegex([]string{"bittrex"}, alphaNumeric("32"), true),
		Keywords:    []string{"bittrex"},
	}

	// validate
	tps := []string{
		generateSampleSecret("bittrex", utils.NewSecret(alphaNumeric("32"))),
	}
	return validate(r, tps, nil)
}

func BittrexSecretKey() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Detected a Bittrex Secret Key, potentially compromising cryptocurrency transactions and financial security.",
		RuleID:      "bittrex-secret-key",
		Regex:       generateSemiGenericRegex([]string{"bittrex"}, alphaNumeric("32"), true),

		Keywords: []string{"bittrex"},
	}

	// validate
	tps := []string{
		generateSampleSecret("bittrex", utils.NewSecret(alphaNumeric("32"))),
	}
	return validate(r, tps, nil)
}
