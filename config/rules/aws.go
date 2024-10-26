package rules

import (
	"github.com/codevault-llc/humblebrag-api/pkg/types"
	regexp "github.com/wasilibs/go-re2"
)

func AWS() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Identified a pattern that may indicate AWS credentials, risking unauthorized cloud resource access and data breaches on AWS platforms.",
		RuleID:      "aws-access-token",
		Regex: regexp.MustCompile(
			"(?:A3T[A-Z0-9]|AKIA|ASIA|ABIA|ACCA)[A-Z0-9]{16}"),
		Keywords: []string{
			"AKIA",
			"ASIA",
			"ABIA",
			"ACCA",
		},
	}

	// validate
	tps := []string{generateSampleSecret("AWS", "AKIALALEMEL33243OLIB")} // gitleaks:allow
	return validate(r, tps, nil)
}
