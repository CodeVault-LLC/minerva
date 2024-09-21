package rules

import (
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
	"github.com/codevault-llc/humblebrag-api/types"
)

func FlickrAccessToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "flickr-access-token",
		Description: "Discovered a Flickr Access Token, posing a risk of unauthorized photo management and potential data leakage.",
		Regex:       generateSemiGenericRegex([]string{"flickr"}, alphaNumeric("32"), true),

		Keywords: []string{
			"flickr",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("flickr", utils.NewSecret(alphaNumeric("32"))),
	}
	return validate(r, tps, nil)
}
