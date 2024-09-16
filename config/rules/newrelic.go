package rules

import (
	"github.com/codevault-llc/humblebrag-api/types"
	"github.com/codevault-llc/humblebrag-api/utils"
)

func NewRelicUserID() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "new-relic-user-api-key",
		Description: "Discovered a New Relic user API Key, which could lead to compromised application insights and performance monitoring.",
		Regex: generateSemiGenericRegex([]string{
			"new-relic",
			"newrelic",
			"new_relic",
		}, `NRAK-[a-z0-9]{27}`, true),

		Keywords: []string{
			"NRAK",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("new-relic", "NRAK-"+utils.NewSecret(alphaNumeric("27"))),
	}
	return validate(r, tps, nil)
}

func NewRelicUserKey() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "new-relic-user-api-id",
		Description: "Found a New Relic user API ID, posing a risk to application monitoring services and data integrity.",
		Regex: generateSemiGenericRegex([]string{
			"new-relic",
			"newrelic",
			"new_relic",
		}, alphaNumeric("64"), true),

		Keywords: []string{
			"new-relic",
			"newrelic",
			"new_relic",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("new-relic", utils.NewSecret(alphaNumeric("64"))),
	}
	return validate(r, tps, nil)
}

func NewRelicBrowserAPIKey() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "new-relic-browser-api-token",
		Description: "Identified a New Relic ingest browser API token, risking unauthorized access to application performance data and analytics.",
		Regex: generateSemiGenericRegex([]string{
			"new-relic",
			"newrelic",
			"new_relic",
		}, `NRJS-[a-f0-9]{19}`, true),

		Keywords: []string{
			"NRJS-",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("new-relic", "NRJS-"+utils.NewSecret(hex("19"))),
	}
	return validate(r, tps, nil)
}

func NewRelicInsertKey() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "new-relic-insert-key",
		Description: "Discovered a New Relic insight insert key, compromising data injection into the platform.",
		Regex: generateSemiGenericRegex([]string{
			"new-relic",
			"newrelic",
			"new_relic",
		}, `NRII-[a-z0-9-]{32}`, true),

		Keywords: []string{
			"NRII-",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("new-relic", "NRII-"+utils.NewSecret(hex("32"))),
	}
	return validate(r, tps, nil)
}
