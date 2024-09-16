package rules

import (
	"github.com/codevault-llc/humblebrag-api/types"
	"github.com/codevault-llc/humblebrag-api/utils"
)

func PulumiAPIToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "pulumi-api-token",
		Description: "Found a Pulumi API token, posing a risk to infrastructure as code services and cloud resource management.",
		Regex:       generateUniqueTokenRegex(`pul-[a-f0-9]{40}`, true),

		Keywords: []string{
			"pul-",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("pulumi-api-token", "pul-"+utils.NewSecret(hex("40"))),
	}
	return validate(r, tps, nil)
}
