package rules

import (
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
	"github.com/codevault-llc/humblebrag-api/types"
)

func MapBox() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Detected a MapBox API token, posing a risk to geospatial services and sensitive location data exposure.",
		RuleID:      "mapbox-api-token",
		Regex:       generateSemiGenericRegex([]string{"mapbox"}, `pk\.[a-z0-9]{60}\.[a-z0-9]{22}`, true),

		Keywords: []string{"mapbox"},
	}

	// validate
	tps := []string{
		generateSampleSecret("mapbox", "pk."+utils.NewSecret(alphaNumeric("60"))+"."+utils.NewSecret(alphaNumeric("22"))),
	}
	return validate(r, tps, nil)
}
