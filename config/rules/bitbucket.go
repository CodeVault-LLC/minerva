package rules

import (
	"github.com/codevault-llc/humblebrag-api/pkg/types"
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
)

func BitBucketClientID() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Discovered a potential Bitbucket Client ID, risking unauthorized repository access and potential codebase exposure.",
		RuleID:      "bitbucket-client-id",
		Regex:       generateSemiGenericRegex([]string{"bitbucket"}, alphaNumeric("32"), true),
		Keywords:    []string{"bitbucket"},
	}

	// validate
	tps := []string{
		generateSampleSecret("bitbucket", utils.NewSecret(alphaNumeric("32"))),
	}
	return validate(r, tps, nil)
}

func BitBucketClientSecret() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Discovered a potential Bitbucket Client Secret, posing a risk of compromised code repositories and unauthorized access.",
		RuleID:      "bitbucket-client-secret",
		Regex:       generateSemiGenericRegex([]string{"bitbucket"}, alphaNumericExtended("64"), true),

		Keywords: []string{"bitbucket"},
	}

	// validate
	tps := []string{
		generateSampleSecret("bitbucket", utils.NewSecret(alphaNumeric("64"))),
	}
	return validate(r, tps, nil)
}
