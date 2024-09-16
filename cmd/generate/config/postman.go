package rules

import (
	"github.com/codevault-llc/humblebrag-api/cmd/generate/secrets"
	"github.com/codevault-llc/humblebrag-api/config"
)

func PostManAPI() *config.Rule {
	// define rule
	r := config.Rule{
		RuleID:      "postman-api-token",
		Description: "Uncovered a Postman API token, potentially compromising API testing and development workflows.",
		Regex:       generateUniqueTokenRegex(`PMAK-(?i)[a-f0-9]{24}\-[a-f0-9]{34}`, true),

		Keywords: []string{
			"PMAK-",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("postmanAPItoken", "PMAK-"+secrets.NewSecret(hex("24"))+"-"+secrets.NewSecret(hex("34"))),
	}
	return validate(r, tps, nil)
}
