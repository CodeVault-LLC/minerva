package rules

import (
	"github.com/codevault-llc/minerva/pkg/types"
	"github.com/codevault-llc/minerva/pkg/utils"
)

func ShippoAPIToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "shippo-api-token",
		Description: "Discovered a Shippo API token, potentially compromising shipping services and customer order data.",
		Regex:       generateUniqueTokenRegex(`shippo_(live|test)_[a-f0-9]{40}`, true),

		Keywords: []string{
			"shippo_",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("shippo", "shippo_live_"+utils.NewSecret(hex("40"))),
		generateSampleSecret("shippo", "shippo_test_"+utils.NewSecret(hex("40"))),
	}
	return validate(r, tps, nil)
}
