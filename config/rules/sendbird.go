package rules

import (
	"github.com/codevault-llc/humblebrag-api/pkg/types"
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
)

func SendbirdAccessToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "sendbird-access-token",
		Description: "Uncovered a Sendbird Access Token, potentially risking unauthorized access to communication services and user data.",
		Regex:       generateSemiGenericRegex([]string{"sendbird"}, hex("40"), true),

		Keywords: []string{
			"sendbird",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("sendbird", utils.NewSecret(hex("40"))),
	}
	return validate(r, tps, nil)
}

func SendbirdAccessID() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "sendbird-access-id",
		Description: "Discovered a Sendbird Access ID, which could compromise chat and messaging platform integrations.",
		Regex:       generateSemiGenericRegex([]string{"sendbird"}, hex8_4_4_4_12(), true),

		Keywords: []string{
			"sendbird",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("sendbird", utils.NewSecret(hex8_4_4_4_12())),
	}
	return validate(r, tps, nil)
}
