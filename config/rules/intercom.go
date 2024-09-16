package rules

import (
	"github.com/codevault-llc/humblebrag-api/types"
	"github.com/codevault-llc/humblebrag-api/utils"
)

func Intercom() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Identified an Intercom API Token, which could compromise customer communication channels and data privacy.",
		RuleID:      "intercom-api-key",
		Regex:       generateSemiGenericRegex([]string{"intercom"}, alphaNumericExtended("60"), true),

		Keywords: []string{"intercom"},
	}

	// validate
	tps := []string{
		generateSampleSecret("intercom", utils.NewSecret(alphaNumericExtended("60"))),
	}
	return validate(r, tps, nil)
}
