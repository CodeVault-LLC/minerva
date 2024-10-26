package rules

import (
	"github.com/codevault-llc/humblebrag-api/pkg/types"
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
)

func KucoinAccessToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "kucoin-access-token",
		Description: "Found a Kucoin Access Token, risking unauthorized access to cryptocurrency exchange services and transactions.",
		Regex:       generateSemiGenericRegex([]string{"kucoin"}, hex("24"), true),

		Keywords: []string{
			"kucoin",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("kucoin", utils.NewSecret(hex("24"))),
	}
	return validate(r, tps, nil)
}

func KucoinSecretKey() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "kucoin-secret-key",
		Description: "Discovered a Kucoin Secret Key, which could lead to compromised cryptocurrency operations and financial data breaches.",
		Regex:       generateSemiGenericRegex([]string{"kucoin"}, hex8_4_4_4_12(), true),

		Keywords: []string{
			"kucoin",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("kucoin", utils.NewSecret(hex8_4_4_4_12())),
	}
	return validate(r, tps, nil)
}
