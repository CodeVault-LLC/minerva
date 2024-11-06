package rules

import (
	regexp "github.com/wasilibs/go-re2"

	"github.com/codevault-llc/minerva/pkg/types"
	"github.com/codevault-llc/minerva/pkg/utils"
)

func Duffel() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "duffel-api-token",
		Description: "Uncovered a Duffel API token, which may compromise travel platform integrations and sensitive customer data.",
		Regex:       regexp.MustCompile(`duffel_(test|live)_(?i)[a-z0-9_\-=]{43}`),
		Keywords:    []string{"duffel"},
	}

	// validate
	tps := []string{
		generateSampleSecret("duffel", "duffel_test_"+utils.NewSecret(alphaNumericExtended("43"))),
	}
	return validate(r, tps, nil)
}
