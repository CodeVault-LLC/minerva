package rules

import (
	"github.com/codevault-llc/minerva/pkg/types"
	"github.com/codevault-llc/minerva/pkg/utils"
)

func Intra42ClientSecret() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Found a Intra42 client secret, which could lead to unauthorized access to the 42School API and sensitive data.",
		RuleID:      "intra42-client-secret",
		Regex:       generateUniqueTokenRegex(`s-s4t2(?:ud|af)-[abcdef0123456789]{64}`, true),
		Keywords: []string{
			"intra",
			"s-s4t2ud-",
			"s-s4t2af-",
		},
	}

	// validate
	tps := []string{
		"clientSecret := \"s-s4t2ud-" + utils.NewSecret(hex("64")) + "\"",
		"clientSecret := \"s-s4t2af-" + utils.NewSecret(hex("64")) + "\"",
		"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", // gitleaks:allow
		"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", // gitleaks:allow
	}
	return validate(r, tps, nil)
}
