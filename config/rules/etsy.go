package rules

import (
	"github.com/codevault-llc/minerva/pkg/types"
	"github.com/codevault-llc/minerva/pkg/utils"
)

func EtsyAccessToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "etsy-access-token",
		Description: "Found an Etsy Access Token, potentially compromising Etsy shop management and customer data.",
		Regex:       generateSemiGenericRegex([]string{"etsy"}, alphaNumeric("24"), true),

		Keywords: []string{
			"etsy",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("etsy", utils.NewSecret(alphaNumeric("24"))),
	}
	return validate(r, tps, nil)
}
