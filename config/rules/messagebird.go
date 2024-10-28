package rules

import (
	"github.com/codevault-llc/humblebrag-api/pkg/types"
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
)

func MessageBirdAPIToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Found a MessageBird API token, risking unauthorized access to communication platforms and message data.",
		RuleID:      "messagebird-api-token",
		Regex: generateSemiGenericRegex([]string{
			"messagebird",
			"message-bird",
			"message_bird",
		}, alphaNumeric("25"), true),

		Keywords: []string{
			"messagebird",
			"message-bird",
			"message_bird",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("messagebird", utils.NewSecret(alphaNumeric("25"))),
		generateSampleSecret("message-bird", utils.NewSecret(alphaNumeric("25"))),
		generateSampleSecret("message_bird", utils.NewSecret(alphaNumeric("25"))),
	}
	return validate(r, tps, nil)
}

func MessageBirdClientID() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Discovered a MessageBird client ID, potentially compromising API integrations and sensitive communication data.",
		RuleID:      "messagebird-client-id",
		Regex: generateSemiGenericRegex([]string{
			"messagebird",
			"message-bird",
			"message_bird",
		}, hex8_4_4_4_12(), true),

		Keywords: []string{
			"messagebird",
			"message-bird",
			"message_bird",
		},
	}

	// validate
	tps := []string{
		`const MessageBirdClientID = "12345678-ABCD-ABCD-ABCD-1234567890AB"`, // gitleaks:allow
	}
	return validate(r, tps, nil)
}
