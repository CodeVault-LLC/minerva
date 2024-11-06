package rules

import (
	"github.com/codevault-llc/minerva/pkg/types"
	"github.com/codevault-llc/minerva/pkg/utils"
)

func PostManAPI() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "postman-api-token",
		Description: "Uncovered a Postman API token, potentially compromising API testing and development workflows.",
		Regex:       generateUniqueTokenRegex(`PMAK-(?i)[a-f0-9]{24}\-[a-f0-9]{34}`, true),

		Keywords: []string{
			"PMAK-",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("postmanAPItoken", "PMAK-"+utils.NewSecret(hex("24"))+"-"+utils.NewSecret(hex("34"))),
	}
	return validate(r, tps, nil)
}
