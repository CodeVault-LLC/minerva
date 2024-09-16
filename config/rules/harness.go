package rules

import (
	regexp "github.com/wasilibs/go-re2"

	"github.com/codevault-llc/humblebrag-api/types"
	"github.com/codevault-llc/humblebrag-api/utils"
)

func HarnessApiKey() *types.Rule {
	// Define types.Rule for Harness Personal Access Token (PAT) and Service Account Token (SAT)
	r := types.Rule{
		Description: "Identified a Harness Access Token (PAT or SAT), risking unauthorized access to a Harness account.",
		RuleID:      "harness-api-key",
		Regex:       regexp.MustCompile(`((?:pat|sat)\.[a-zA-Z0-9]{22}\.[a-zA-Z0-9]{24}\.[a-zA-Z0-9]{20})`),
		Keywords:    []string{"pat.", "sat."},
	}

	// Generate a sample secret for validation
	tps := []string{
		generateSampleSecret("harness", "pat."+utils.NewSecret(alphaNumeric("22"))+"."+utils.NewSecret(alphaNumeric("24"))+"."+utils.NewSecret(alphaNumeric("20"))),
		generateSampleSecret("harness", "sat."+utils.NewSecret(alphaNumeric("22"))+"."+utils.NewSecret(alphaNumeric("24"))+"."+utils.NewSecret(alphaNumeric("20"))),
	}

	// Validate the types.Rule
	return validate(r, tps, nil)
}
