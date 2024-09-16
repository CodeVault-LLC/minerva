package rules

import (
	regexp "github.com/wasilibs/go-re2"

	"github.com/codevault-llc/humblebrag-api/cmd/generate/secrets"
	"github.com/codevault-llc/humblebrag-api/config"
)

func Clojars() *config.Rule {
	// define rule
	r := config.Rule{
		Description: "Uncovered a possible Clojars API token, risking unauthorized access to Clojure libraries and potential code manipulation.",
		RuleID:      "clojars-api-token",
		Regex:       regexp.MustCompile(`(?i)(CLOJARS_)[a-z0-9]{60}`),
		Keywords:    []string{"clojars"},
	}

	// validate
	tps := []string{
		generateSampleSecret("clojars", "CLOJARS_"+secrets.NewSecret(alphaNumeric("60"))),
	}
	return validate(r, tps, nil)
}
