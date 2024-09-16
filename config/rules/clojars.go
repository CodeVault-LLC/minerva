package rules

import (
	regexp "github.com/wasilibs/go-re2"

	"github.com/codevault-llc/humblebrag-api/types"
	"github.com/codevault-llc/humblebrag-api/utils"
)

func Clojars() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Uncovered a possible Clojars API token, risking unauthorized access to Clojure libraries and potential code manipulation.",
		RuleID:      "clojars-api-token",
		Regex:       regexp.MustCompile(`(?i)(CLOJARS_)[a-z0-9]{60}`),
		Keywords:    []string{"clojars"},
	}

	// validate
	tps := []string{
		generateSampleSecret("clojars", "CLOJARS_"+utils.NewSecret(alphaNumeric("60"))),
	}
	return validate(r, tps, nil)
}
