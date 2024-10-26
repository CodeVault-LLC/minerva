package rules

import (
	"github.com/codevault-llc/humblebrag-api/pkg/types"
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
)

func TwitterAPIKey() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Identified a Twitter API Key, which may compromise Twitter application integrations and user data security.",
		RuleID:      "twitter-api-key",
		Regex:       generateSemiGenericRegex([]string{"twitter"}, alphaNumeric("25"), true),
		Keywords:    []string{"twitter"},
	}

	// validate
	tps := []string{
		generateSampleSecret("twitter", utils.NewSecret(alphaNumeric("25"))),
	}
	return validate(r, tps, nil)
}

func TwitterAPISecret() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Found a Twitter API Secret, risking the security of Twitter app integrations and sensitive data access.",
		RuleID:      "twitter-api-secret",
		Regex:       generateSemiGenericRegex([]string{"twitter"}, alphaNumeric("50"), true),
		Keywords:    []string{"twitter"},
	}

	// validate
	tps := []string{
		generateSampleSecret("twitter", utils.NewSecret(alphaNumeric("50"))),
	}
	return validate(r, tps, nil)
}

func TwitterBearerToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Discovered a Twitter Bearer Token, potentially compromising API access and data retrieval from Twitter.",
		RuleID:      "twitter-bearer-token",
		Regex:       generateSemiGenericRegex([]string{"twitter"}, "A{22}[a-zA-Z0-9%]{80,100}", true),

		Keywords: []string{"twitter"},
	}

	// validate
	tps := []string{
		generateSampleSecret("twitter", utils.NewSecret("A{22}[a-zA-Z0-9%]{80,100}")),
	}
	return validate(r, tps, nil)
}

func TwitterAccessToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Detected a Twitter Access Token, posing a risk of unauthorized account operations and social media data exposure.",
		RuleID:      "twitter-access-token",
		Regex:       generateSemiGenericRegex([]string{"twitter"}, "[0-9]{15,25}-[a-zA-Z0-9]{20,40}", true),
		Keywords:    []string{"twitter"},
	}

	// validate
	tps := []string{
		generateSampleSecret("twitter", utils.NewSecret("[0-9]{15,25}-[a-zA-Z0-9]{20,40}")),
	}
	return validate(r, tps, nil)
}

func TwitterAccessSecret() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Uncovered a Twitter Access Secret, potentially risking unauthorized Twitter integrations and data breaches.",
		RuleID:      "twitter-access-secret",
		Regex:       generateSemiGenericRegex([]string{"twitter"}, alphaNumeric("45"), true),
		Keywords:    []string{"twitter"},
	}

	// validate
	tps := []string{
		generateSampleSecret("twitter", utils.NewSecret(alphaNumeric("45"))),
	}
	return validate(r, tps, nil)
}
