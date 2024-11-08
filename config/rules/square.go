package rules

import (
	"github.com/codevault-llc/minerva/pkg/types"
	"github.com/codevault-llc/minerva/pkg/utils"
)

func SquareAccessToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "square-access-token",
		Description: "Detected a Square Access Token, risking unauthorized payment processing and financial transaction exposure.",
		Regex:       generateUniqueTokenRegex(`(EAAA|sq0atp-)[0-9A-Za-z\-_]{22,60}`, true),
		Keywords:    []string{"sq0atp-", "EAAA"},
	}

	// validate
	tps := []string{
		generateSampleSecret("square", utils.NewSecret(`sq0atp-[0-9A-Za-z\-_]{22}`)),
		"ARG token=sq0atp-812erere3wewew45678901",                                    // gitleaks:allow
		"ARG token=EAAAlsBxkkVgvmr7FasTFbM6VUGZ31EJ4jZKTJZySgElBDJ_wyafHuBFquFexY7E", // gitleaks:allow",
	}
	return validate(r, tps, nil)
}

func SquareSecret() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "square-secret",
		Description: "Square Secret",
		Regex:       generateUniqueTokenRegex(`sq0csp-[0-9A-Za-z\\-_]{43}`, true),
		Keywords:    []string{"sq0csp-"},
	}

	// validate
	tps := []string{
		generateSampleSecret("square", utils.NewSecret(`sq0csp-[0-9A-Za-z\\-_]{43}`)),
		`value: "sq0csp-0p9h7g6f4s3s3s3-4a3ardgwa6ADRDJDDKUFYDYDYDY"`, // gitleaks:allow
	}
	return validate(r, tps, nil)
}
