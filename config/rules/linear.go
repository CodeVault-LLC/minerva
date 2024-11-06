package rules

import (
	regexp "github.com/wasilibs/go-re2"

	"github.com/codevault-llc/minerva/pkg/types"
	"github.com/codevault-llc/minerva/pkg/utils"
)

func LinearAPIToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Detected a Linear API Token, posing a risk to project management tools and sensitive task data.",
		RuleID:      "linear-api-key",
		Regex:       regexp.MustCompile(`lin_api_(?i)[a-z0-9]{40}`),
		Keywords:    []string{"lin_api_"},
	}

	// validate
	tps := []string{
		generateSampleSecret("linear", "lin_api_"+utils.NewSecret(alphaNumeric("40"))),
	}
	return validate(r, tps, nil)
}

func LinearClientSecret() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Identified a Linear Client Secret, which may compromise secure integrations and sensitive project management data.",
		RuleID:      "linear-client-secret",
		Regex:       generateSemiGenericRegex([]string{"linear"}, hex("32"), true),
		Keywords:    []string{"linear"},
	}

	// validate
	tps := []string{
		generateSampleSecret("linear", utils.NewSecret(hex("32"))),
	}
	return validate(r, tps, nil)
}
