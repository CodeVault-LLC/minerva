package rules

import (
	"github.com/codevault-llc/humblebrag-api/cmd/generate/secrets"
	"github.com/codevault-llc/humblebrag-api/config"
)

func Contentful() *config.Rule {
	// define rule
	r := config.Rule{
		Description: "Discovered a Contentful delivery API token, posing a risk to content management systems and data integrity.",
		RuleID:      "contentful-delivery-api-token",
		Regex: generateSemiGenericRegex([]string{"contentful"},
			alphaNumericExtended("43"), true),
		Keywords: []string{"contentful"},
	}

	// validate
	tps := []string{
		generateSampleSecret("contentful", secrets.NewSecret(alphaNumeric("43"))),
	}
	return validate(r, tps, nil)
}
