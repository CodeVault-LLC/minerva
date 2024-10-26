package rules

import (
	"github.com/codevault-llc/humblebrag-api/pkg/types"
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
)

func AdafruitAPIKey() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Identified a potential Adafruit API Key, which could lead to unauthorized access to Adafruit services and sensitive data exposure.",
		RuleID:      "adafruit-api-key",
		Regex:       generateSemiGenericRegex([]string{"adafruit"}, alphaNumericExtendedShort("32"), true),
		Keywords:    []string{"adafruit"},
	}

	// validate
	tps := []string{
		generateSampleSecret("adafruit", utils.NewSecret(alphaNumericExtendedShort("32"))),
	}
	return validate(r, tps, nil)
}
