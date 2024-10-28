package rules

import (
	"github.com/codevault-llc/humblebrag-api/pkg/types"
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
)

func ConfluentSecretKey() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "confluent-secret-key",
		Description: "Found a Confluent Secret Key, potentially risking unauthorized operations and data access within Confluent services.",
		Regex:       generateSemiGenericRegex([]string{"confluent"}, alphaNumeric("64"), true),
		Keywords: []string{
			"confluent",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("confluent", utils.NewSecret(alphaNumeric("64"))),
	}
	return validate(r, tps, nil)
}

func ConfluentAccessToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "confluent-access-token",
		Description: "Identified a Confluent Access Token, which could compromise access to streaming data platforms and sensitive data flow.",
		Regex:       generateSemiGenericRegex([]string{"confluent"}, alphaNumeric("16"), true),

		Keywords: []string{
			"confluent",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("confluent", utils.NewSecret(alphaNumeric("16"))),
	}
	return validate(r, tps, nil)
}
