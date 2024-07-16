package rules

import (
	"github.com/codevault-llc/humblebrag-api/cmd/generate/secrets"
	"github.com/codevault-llc/humblebrag-api/config"
)

func ZendeskSecretKey() *config.Rule {
	// define rule
	r := config.Rule{
		RuleID:      "zendesk-secret-key",
		Description: "Detected a Zendesk Secret Key, risking unauthorized access to customer support services and sensitive ticketing data.",
		Regex:       generateSemiGenericRegex([]string{"zendesk"}, alphaNumeric("40"), true),
		Keywords: []string{
			"zendesk",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("zendesk", secrets.NewSecret(alphaNumeric("40"))),
	}
	return validate(r, tps, nil)
}
