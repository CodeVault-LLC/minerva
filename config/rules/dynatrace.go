package rules

import (
	regexp "github.com/wasilibs/go-re2"

	"github.com/codevault-llc/minerva/pkg/types"
	"github.com/codevault-llc/minerva/pkg/utils"
)

func Dynatrace() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Detected a Dynatrace API token, potentially risking application performance monitoring and data exposure.",
		RuleID:      "dynatrace-api-token",
		Regex:       regexp.MustCompile(`dt0c01\.(?i)[a-z0-9]{24}\.[a-z0-9]{64}`),
		Keywords:    []string{"dynatrace"},
	}

	// validate
	tps := []string{
		generateSampleSecret("dynatrace", "dt0c01."+utils.NewSecret(alphaNumeric("24"))+"."+utils.NewSecret(alphaNumeric("64"))),
	}
	return validate(r, tps, nil)
}
