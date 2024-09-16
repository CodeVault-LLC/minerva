package rules

import (
	"github.com/codevault-llc/humblebrag-api/types"
	"github.com/codevault-llc/humblebrag-api/utils"
)

func CodecovAccessToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "codecov-access-token",
		Description: "Found a pattern resembling a Codecov Access Token, posing a risk of unauthorized access to code coverage reports and sensitive data.",
		Regex:       generateSemiGenericRegex([]string{"codecov"}, alphaNumeric("32"), true),
		Keywords: []string{
			"codecov",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("codecov", utils.NewSecret(alphaNumeric("32"))),
	}
	return validate(r, tps, nil)
}
