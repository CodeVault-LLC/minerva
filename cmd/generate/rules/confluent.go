package rules

import (
	"github.com/codevault-llc/humblebrag-api/cmd/generate/secrets"
	"github.com/codevault-llc/humblebrag-api/config"
)

func ConfluentSecretKey() *config.Rule {
	// define rule
	r := config.Rule{
		RuleID:      "confluent-secret-key",
		Description: "Found a Confluent Secret Key, potentially risking unauthorized operations and data access within Confluent services.",
		Regex:       generateSemiGenericRegex([]string{"confluent"}, alphaNumeric("64"), true),
		Keywords: []string{
			"confluent",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("confluent", secrets.NewSecret(alphaNumeric("64"))),
	}
	return validate(r, tps, nil)
}

func ConfluentAccessToken() *config.Rule {
	// define rule
	r := config.Rule{
		RuleID:      "confluent-access-token",
		Description: "Identified a Confluent Access Token, which could compromise access to streaming data platforms and sensitive data flow.",
		Regex:       generateSemiGenericRegex([]string{"confluent"}, alphaNumeric("16"), true),

		Keywords: []string{
			"confluent",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("confluent", secrets.NewSecret(alphaNumeric("16"))),
	}
	return validate(r, tps, nil)
}
