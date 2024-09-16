package rules

import (
	regexp "github.com/wasilibs/go-re2"

	"github.com/codevault-llc/humblebrag-api/types"
	"github.com/codevault-llc/humblebrag-api/utils"
)

func FrameIO() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Found a Frame.io API token, potentially compromising video collaboration and project management.",
		RuleID:      "frameio-api-token",
		Regex:       regexp.MustCompile(`fio-u-(?i)[a-z0-9\-_=]{64}`),
		Keywords:    []string{"fio-u-"},
	}

	// validate
	tps := []string{
		generateSampleSecret("frameio", "fio-u-"+utils.NewSecret(alphaNumericExtended("64"))),
	}
	return validate(r, tps, nil)
}
