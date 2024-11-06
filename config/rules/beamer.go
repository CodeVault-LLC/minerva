package rules

import (
	"github.com/codevault-llc/minerva/pkg/types"
	"github.com/codevault-llc/minerva/pkg/utils"
)

func Beamer() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Detected a Beamer API token, potentially compromising content management and exposing sensitive notifications and updates.",
		RuleID:      "beamer-api-token",
		Regex: generateSemiGenericRegex([]string{"beamer"},
			`b_[a-z0-9=_\-]{44}`, true),
		Keywords: []string{"beamer"},
	}

	// validate
	tps := []string{
		generateSampleSecret("beamer", "b_"+utils.NewSecret(alphaNumericExtended("44"))),
	}
	return validate(r, tps, nil)
}
