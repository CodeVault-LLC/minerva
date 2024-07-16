package rules

import (
	"regexp"

	"github.com/codevault-llc/humblebrag-api/cmd/generate/secrets"
	"github.com/codevault-llc/humblebrag-api/config"
)

func PyPiUploadToken() *config.Rule {
	// define rule
	r := config.Rule{
		Description: "Discovered a PyPI upload token, potentially compromising Python package distribution and repository integrity.",
		RuleID:      "pypi-upload-token",
		Regex: regexp.MustCompile(
			`pypi-AgEIcHlwaS5vcmc[A-Za-z0-9\-_]{50,1000}`),
		Keywords: []string{
			"pypi-AgEIcHlwaS5vcmc",
		},
	}

	// validate
	tps := []string{"pypiToken := \"pypi-AgEIcHlwaS5vcmc" + secrets.NewSecret(hex("32")) +
		secrets.NewSecret(hex("32")) + "\""}
	return validate(r, tps, nil)
}
