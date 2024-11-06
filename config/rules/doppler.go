package rules

import (
	regexp "github.com/wasilibs/go-re2"

	"github.com/codevault-llc/minerva/pkg/types"
	"github.com/codevault-llc/minerva/pkg/utils"
)

func Doppler() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Discovered a Doppler API token, posing a risk to environment and secrets management security.",
		RuleID:      "doppler-api-token",
		Regex:       regexp.MustCompile(`(dp\.pt\.)(?i)[a-z0-9]{43}`),
		Keywords:    []string{"doppler"},
	}

	// validate
	tps := []string{
		generateSampleSecret("doppler", "dp.pt."+utils.NewSecret(alphaNumeric("43"))),
	}
	return validate(r, tps, nil)
}

// TODO add additional doppler formats:
// https://docs.doppler.com/reference/auth-token-formats
