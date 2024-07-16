package rules

import (
	"github.com/codevault-llc/humblebrag-api/cmd/generate/secrets"
	"github.com/codevault-llc/humblebrag-api/config"
)

func Airtable() *config.Rule {
	// define rule
	r := config.Rule{
		Description: "Uncovered a possible Airtable API Key, potentially compromising database access and leading to data leakage or alteration.",
		RuleID:      "airtable-api-key",
		Regex:       generateSemiGenericRegex([]string{"airtable"}, alphaNumeric("17"), true),
		Keywords:    []string{"airtable"},
	}

	// validate
	tps := []string{
		generateSampleSecret("airtable", secrets.NewSecret(alphaNumeric("17"))),
	}
	return validate(r, tps, nil)
}
