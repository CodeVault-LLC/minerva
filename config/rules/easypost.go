package rules

import (
	regexp "github.com/wasilibs/go-re2"

	"github.com/codevault-llc/humblebrag-api/types"
	"github.com/codevault-llc/humblebrag-api/utils"
)

func EasyPost() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Identified an EasyPost API token, which could lead to unauthorized postal and shipment service access and data exposure.",
		RuleID:      "easypost-api-token",
		Regex:       regexp.MustCompile(`\bEZAK(?i)[a-z0-9]{54}`),
		Keywords:    []string{"EZAK"},
	}

	// validate
	tps := []string{
		generateSampleSecret("EZAK", "EZAK"+utils.NewSecret(alphaNumeric("54"))),
	}
	return validate(r, tps, nil)
}

func EasyPostTestAPI() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Detected an EasyPost test API token, risking exposure of test environments and potentially sensitive shipment data.",
		RuleID:      "easypost-test-api-token",
		Regex:       regexp.MustCompile(`\bEZTK(?i)[a-z0-9]{54}`),
		Keywords:    []string{"EZTK"},
	}

	// validate
	tps := []string{
		generateSampleSecret("EZTK", "EZTK"+utils.NewSecret(alphaNumeric("54"))),
	}
	return validate(r, tps, nil)
}
