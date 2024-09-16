package rules

import (
	"github.com/codevault-llc/humblebrag-api/types"
	"github.com/codevault-llc/humblebrag-api/utils"
)

func Airtable() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Uncovered a possible Airtable API Key, potentially compromising database access and leading to data leakage or alteration.",
		RuleID:      "airtable-api-key",
		Regex:       generateSemiGenericRegex([]string{"airtable"}, alphaNumeric("17"), true),
		Keywords:    []string{"airtable"},
	}

	// validate
	tps := []string{
		generateSampleSecret("airtable", utils.NewSecret(alphaNumeric("17"))),
	}
	return validate(r, tps, nil)
}
