package rules

import (
	"github.com/codevault-llc/humblebrag-api/pkg/types"
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
)

func AdobeClientID() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Detected a pattern that resembles an Adobe OAuth Web Client ID, posing a risk of compromised Adobe integrations and data breaches.",
		RuleID:      "adobe-client-id",
		Regex:       generateSemiGenericRegex([]string{"adobe"}, hex("32"), true),
		Keywords:    []string{"adobe"},
	}

	// validate
	tps := []string{
		generateSampleSecret("adobe", utils.NewSecret(hex("32"))),
	}
	return validate(r, tps, nil)
}

func AdobeClientSecret() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Discovered a potential Adobe Client Secret, which, if exposed, could allow unauthorized Adobe service access and data manipulation.",
		RuleID:      "adobe-client-secret",
		Regex:       generateUniqueTokenRegex(`(p8e-)(?i)[a-z0-9]{32}`, true),
		Keywords:    []string{"p8e-"},
	}

	// validate
	tps := []string{
		"adobeClient := \"p8e-" + utils.NewSecret(hex("32")) + "\"",
	}
	return validate(r, tps, nil)
}
