package rules

import (
	"github.com/codevault-llc/humblebrag-api/pkg/types"
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
)

func SendInBlueAPIToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "sendinblue-api-token",
		Description: "Identified a Sendinblue API token, which may compromise email marketing services and subscriber data privacy.",
		Regex:       generateUniqueTokenRegex(`xkeysib-[a-f0-9]{64}\-(?i)[a-z0-9]{16}`, true),

		Keywords: []string{
			"xkeysib-",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("sendinblue", "xkeysib-"+utils.NewSecret(hex("64"))+"-"+utils.NewSecret(alphaNumeric("16"))),
	}
	return validate(r, tps, nil)
}
