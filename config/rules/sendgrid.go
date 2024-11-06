package rules

import (
	"github.com/codevault-llc/minerva/pkg/types"
	"github.com/codevault-llc/minerva/pkg/utils"
)

func SendGridAPIToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "sendgrid-api-token",
		Description: "Detected a SendGrid API token, posing a risk of unauthorized email service operations and data exposure.",
		Regex:       generateUniqueTokenRegex(`SG\.(?i)[a-z0-9=_\-\.]{66}`, true),

		Keywords: []string{
			"SG.",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("sengridAPIToken", "SG."+utils.NewSecret(alphaNumericExtended("66"))),
	}
	return validate(r, tps, nil)
}
