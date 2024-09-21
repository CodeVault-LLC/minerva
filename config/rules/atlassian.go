package rules

import (
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
	"github.com/codevault-llc/humblebrag-api/types"
)

func Atlassian() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Detected an Atlassian API token, posing a threat to project management and collaboration tool security and data confidentiality.",
		RuleID:      "atlassian-api-token",
		Regex: generateSemiGenericRegex([]string{
			"atlassian", "confluence", "jira"}, alphaNumeric("24"), true),
		Keywords: []string{"atlassian", "confluence", "jira"},
	}

	// validate
	tps := []string{
		generateSampleSecret("atlassian", utils.NewSecret(alphaNumeric("24"))),
		generateSampleSecret("confluence", utils.NewSecret(alphaNumeric("24"))),
	}
	return validate(r, tps, nil)
}
