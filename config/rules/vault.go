package rules

import (
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
	"github.com/codevault-llc/humblebrag-api/types"
)

func VaultServiceToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Identified a Vault Service Token, potentially compromising infrastructure security and access to sensitive credentials.",
		RuleID:      "vault-service-token",
		Regex:       generateUniqueTokenRegex(`hvs\.[a-z0-9_-]{90,100}`, true),
		Keywords:    []string{"hvs"},
	}

	// validate
	tps := []string{
		generateSampleSecret("vault", "hvs."+utils.NewSecret(alphaNumericExtendedShort("90"))),
	}
	return validate(r, tps, nil)
}

func VaultBatchToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Detected a Vault Batch Token, risking unauthorized access to secret management services and sensitive data.",
		RuleID:      "vault-batch-token",
		Regex:       generateUniqueTokenRegex(`hvb\.[a-z0-9_-]{138,212}`, true),
		Keywords:    []string{"hvb"},
	}

	// validate
	tps := []string{
		generateSampleSecret("vault", "hvb."+utils.NewSecret(alphaNumericExtendedShort("138"))),
	}
	return validate(r, tps, nil)
}
