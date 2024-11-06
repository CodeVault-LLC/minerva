package rules

import (
	"github.com/codevault-llc/minerva/pkg/types"
	"github.com/codevault-llc/minerva/pkg/utils"
)

func ZendeskSecretKey() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "zendesk-secret-key",
		Description: "Detected a Zendesk Secret Key, risking unauthorized access to customer support services and sensitive ticketing data.",
		Regex:       generateSemiGenericRegex([]string{"zendesk"}, alphaNumeric("40"), true),
		Keywords: []string{
			"zendesk",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("zendesk", utils.NewSecret(alphaNumeric("40"))),
	}
	return validate(r, tps, nil)
}
