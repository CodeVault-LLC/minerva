package rules

import (
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
	"github.com/codevault-llc/humblebrag-api/types"
)

func GrafanaApiKey() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Identified a Grafana API key, which could compromise monitoring dashboards and sensitive data analytics.",
		RuleID:      "grafana-api-key",

		Regex:    generateUniqueTokenRegex(`eyJrIjoi[A-Za-z0-9]{70,400}={0,2}`, true),
		Keywords: []string{"eyJrIjoi"},
	}

	// validate
	tps := []string{
		generateSampleSecret("grafana-api-key",
			"eyJrIjoi"+
				utils.NewSecret(alphaNumeric("70"))),
	}
	return validate(r, tps, nil)
}

func GrafanaCloudApiToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Found a Grafana cloud API token, risking unauthorized access to cloud-based monitoring services and data exposure.",
		RuleID:      "grafana-cloud-api-token",

		Regex:    generateUniqueTokenRegex(`glc_[A-Za-z0-9+/]{32,400}={0,2}`, true),
		Keywords: []string{"glc_"},
	}

	// validate
	tps := []string{
		generateSampleSecret("grafana-cloud-api-token",
			"glc_"+
				utils.NewSecret(alphaNumeric("32"))),
	}
	return validate(r, tps, nil)
}

func GrafanaServiceAccountToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Discovered a Grafana service account token, posing a risk of compromised monitoring services and data integrity.",
		RuleID:      "grafana-service-account-token",

		Regex:    generateUniqueTokenRegex(`glsa_[A-Za-z0-9]{32}_[A-Fa-f0-9]{8}`, true),
		Keywords: []string{"glsa_"},
	}

	// validate
	tps := []string{
		generateSampleSecret("grafana-service-account-token",
			"glsa_"+
				utils.NewSecret(alphaNumeric("32"))+
				"_"+
				utils.NewSecret((hex("8")))),
	}
	return validate(r, tps, nil)
}
