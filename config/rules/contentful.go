package rules

import (
	"github.com/codevault-llc/minerva/pkg/types"
	"github.com/codevault-llc/minerva/pkg/utils"
)

func Contentful() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Discovered a Contentful delivery API token, posing a risk to content management systems and data integrity.",
		RuleID:      "contentful-delivery-api-token",
		Regex: generateSemiGenericRegex([]string{"contentful"},
			alphaNumericExtended("43"), true),
		Keywords: []string{"contentful"},
	}

	// validate
	tps := []string{
		generateSampleSecret("contentful", utils.NewSecret(alphaNumeric("43"))),
	}
	return validate(r, tps, nil)
}
