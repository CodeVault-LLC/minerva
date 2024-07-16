package rules

import (
	"regexp"

	"github.com/codevault-llc/humblebrag-api/cmd/generate/secrets"
	"github.com/codevault-llc/humblebrag-api/config"
)

func FrameIO() *config.Rule {
	// define rule
	r := config.Rule{
		Description: "Found a Frame.io API token, potentially compromising video collaboration and project management.",
		RuleID:      "frameio-api-token",
		Regex:       regexp.MustCompile(`fio-u-(?i)[a-z0-9\-_=]{64}`),
		Keywords:    []string{"fio-u-"},
	}

	// validate
	tps := []string{
		generateSampleSecret("frameio", "fio-u-"+secrets.NewSecret(alphaNumericExtended("64"))),
	}
	return validate(r, tps, nil)
}
