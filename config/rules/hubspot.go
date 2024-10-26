package rules

import "github.com/codevault-llc/humblebrag-api/pkg/types"

func HubSpot() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Found a HubSpot API Token, posing a risk to CRM data integrity and unauthorized marketing operations.",
		RuleID:      "hubspot-api-key",
		Regex: generateSemiGenericRegex([]string{"hubspot"},
			`[0-9A-F]{8}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{12}`, true),

		Keywords: []string{"hubspot"},
	}

	// validate
	tps := []string{
		`const hubspotKey = "12345678-ABCD-ABCD-ABCD-1234567890AB"`, // gitleaks:allow
	}
	return validate(r, tps, nil)
}
