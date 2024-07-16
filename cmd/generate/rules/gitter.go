package rules

import (
	"github.com/codevault-llc/humblebrag-api/cmd/generate/secrets"
	"github.com/codevault-llc/humblebrag-api/config"
)

func GitterAccessToken() *config.Rule {
	// define rule
	r := config.Rule{
		RuleID:      "gitter-access-token",
		Description: "Uncovered a Gitter Access Token, which may lead to unauthorized access to chat and communication services.",
		Regex: generateSemiGenericRegex([]string{"gitter"},
			alphaNumericExtendedShort("40"), true),

		Keywords: []string{
			"gitter",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("gitter",
			secrets.NewSecret(alphaNumericExtendedShort("40"))),
	}
	return validate(r, tps, nil)
}
