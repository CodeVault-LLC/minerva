package rules

import (
	"github.com/codevault-llc/minerva/pkg/types"
	"github.com/codevault-llc/minerva/pkg/utils"
)

func ScalingoAPIToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Found a Scalingo API token, posing a risk to cloud platform services and application deployment security.",
		RuleID:      "scalingo-api-token",
		Regex:       generateUniqueTokenRegex(`tk-us-[a-zA-Z0-9-_]{48}`, false),
		Keywords:    []string{"tk-us-"},
	}

	// validate
	tps := []string{
		generateSampleSecret("scalingo", "tk-us-"+utils.NewSecret(alphaNumericExtendedShort("48"))),
		`scalingo_api_token = "tk-us-loys7ib9yrxcys_ta2sq85mjar6lgcsspkd9x61s7h5epf_-"`, // gitleaks:allow
	}
	return validate(r, tps, nil)
}
