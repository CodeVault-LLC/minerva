package rules

import (
	"github.com/codevault-llc/humblebrag-api/types"
	"github.com/codevault-llc/humblebrag-api/utils"
)

func AlgoliaApiKey() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Identified an Algolia API Key, which could result in unauthorized search operations and data exposure on Algolia-managed platforms.",
		RuleID:      "algolia-api-key",
		Regex:       generateSemiGenericRegex([]string{"algolia"}, `[a-z0-9]{32}`, true),
		Keywords:    []string{"algolia"},
	}

	// validate
	tps := []string{
		"algolia_key := " + utils.NewSecret(hex("32")),
	}
	return validate(r, tps, nil)
}