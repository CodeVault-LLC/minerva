package rules

import (
	regexp "github.com/wasilibs/go-re2"

	"github.com/codevault-llc/humblebrag-api/types"
	"github.com/codevault-llc/humblebrag-api/utils"
)

// TODO this one could probably use some work
func GCPServiceAccount() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Google (GCP) Service-account",
		RuleID:      "gcp-service-account",
		Regex:       regexp.MustCompile(`\"type\": \"service_account\"`),
		Keywords:    []string{`\"type\": \"service_account\"`},
	}

	// validate
	tps := []string{
		`"type": "service_account"`,
	}
	return validate(r, tps, nil)
}

func GCPAPIKey() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "gcp-api-key",
		Description: "Uncovered a GCP API key, which could lead to unauthorized access to Google Cloud services and data breaches.",
		Regex:       generateUniqueTokenRegex(`AIza[0-9A-Za-z\\-_]{35}`, true),

		Keywords: []string{
			"AIza",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("gcp", utils.NewSecret(`AIza[0-9A-Za-z\\-_]{35}`)),
	}
	return validate(r, tps, nil)
}

func GCPOAuth() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "gcp-oauth",
		Description: "Uncovered a GCP OAuth key, which could lead to unauthorized access to Google Cloud services and data breaches.",
		Regex:       generateUniqueTokenRegex(`[0-9]+-[0-9A-Za-z_]{32}\.apps\.googleusercontent\.com`, true),

		Keywords: []string{
			"googleusercontent.com",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("gcp", utils.NewSecret(`[0-9]+-[0-9A-Za-z_]{32}\.apps\.googleusercontent\.com`)),
	}
	return validate(r, tps, nil)
}

// Google Tag Manager
func GCPGTM() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "gcp-gtm",
		Description: "Uncovered a Google Tag Manager key, which could lead to unauthorized access to Google Cloud services and data breaches.",
		Regex:       generateUniqueTokenRegex(`GTM-[0-9A-Z]{7}`, true),

		Keywords: []string{
			"GTM",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("gcp", utils.NewSecret(`GTM-[0-9A-Z]{7}`)),
	}
	return validate(r, tps, nil)
}

// Google Analytics
func GCPGA() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "gcp-ga",
		Description: "Uncovered a Google Analytics key, which could lead to unauthorized access to Google Cloud services and data breaches.",
		Regex:       generateUniqueTokenRegex(`UA-[0-9]+-[0-9]+`, true),

		Keywords: []string{
			"UA-",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("gcp", utils.NewSecret(`UA-[0-9]+-[0-9]+`)),
	}
	return validate(r, tps, nil)
}

func GCPGA2() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "gcp-ga2",
		Description: "Uncovered a Google Analytics key, which could lead to unauthorized access to Google Cloud services and data breaches.",
		Regex:       generateUniqueTokenRegex(`G-[0-9A-Z]{10}`, true),

		Keywords: []string{
			"G-",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("gcp", utils.NewSecret(`G-[0-9A-Z]{10}`)),
	}
	return validate(r, tps, nil)
}

func GCPFirebase() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "gcp-firebase",
		Description: "Uncovered a Firebase key, which could lead to unauthorized access to Google Cloud services and data breaches.",
		Regex:       generateUniqueTokenRegex(`[0-9]:[0-9A-Za-z\\-_]{140}`, true),

		Keywords: []string{
			":",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("gcp", utils.NewSecret(`[0-9]:[0-9A-Za-z\\-_]{140}`)),
	}
	return validate(r, tps, nil)
}
