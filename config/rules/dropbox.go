package rules

import (
	"github.com/codevault-llc/humblebrag-api/types"
	"github.com/codevault-llc/humblebrag-api/utils"
)

func DropBoxAPISecret() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Identified a Dropbox API secret, which could lead to unauthorized file access and data breaches in Dropbox storage.",
		RuleID:      "dropbox-api-token",
		Regex:       generateSemiGenericRegex([]string{"dropbox"}, alphaNumeric("15"), true),

		Keywords: []string{"dropbox"},
	}

	// validate
	tps := []string{
		generateSampleSecret("dropbox", utils.NewSecret(alphaNumeric("15"))),
	}
	return validate(r, tps, nil)
}

func DropBoxShortLivedAPIToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "dropbox-short-lived-api-token",
		Description: "Discovered a Dropbox short-lived API token, posing a risk of temporary but potentially harmful data access and manipulation.",
		Regex:       generateSemiGenericRegex([]string{"dropbox"}, `sl\.[a-z0-9\-=_]{135}`, true),
		Keywords:    []string{"dropbox"},
	}

	// validate TODO
	return &r
}

func DropBoxLongLivedAPIToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "dropbox-long-lived-api-token",
		Description: "Found a Dropbox long-lived API token, risking prolonged unauthorized access to cloud storage and sensitive data.",
		Regex:       generateSemiGenericRegex([]string{"dropbox"}, `[a-z0-9]{11}(AAAAAAAAAA)[a-z0-9\-_=]{43}`, true),
		Keywords:    []string{"dropbox"},
	}

	// validate TODO
	return &r
}
