package rules

import (
	"github.com/codevault-llc/minerva/pkg/types"
	"github.com/codevault-llc/minerva/pkg/utils"
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
