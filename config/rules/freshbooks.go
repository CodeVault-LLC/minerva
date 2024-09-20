package rules

import (
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
	"github.com/codevault-llc/humblebrag-api/types"
)

func FreshbooksAccessToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "freshbooks-access-token",
		Description: "Discovered a Freshbooks Access Token, posing a risk to accounting software access and sensitive financial data exposure.",
		Regex:       generateSemiGenericRegex([]string{"freshbooks"}, alphaNumeric("64"), true),

		Keywords: []string{
			"freshbooks",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("freshbooks", utils.NewSecret(alphaNumeric("64"))),
	}
	return validate(r, tps, nil)
}
