package rules

import (
	regexp "github.com/wasilibs/go-re2"

	"github.com/codevault-llc/minerva/pkg/types"
	"github.com/codevault-llc/minerva/pkg/utils"
)

func PyPiUploadToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Discovered a PyPI upload token, potentially compromising Python package distribution and repository integrity.",
		RuleID:      "pypi-upload-token",
		Regex: regexp.MustCompile(
			`pypi-AgEIcHlwaS5vcmc[A-Za-z0-9\-_]{50,1000}`),
		Keywords: []string{
			"pypi-AgEIcHlwaS5vcmc",
		},
	}

	// validate
	tps := []string{"pypiToken := \"pypi-AgEIcHlwaS5vcmc" + utils.NewSecret(hex("32")) +
		utils.NewSecret(hex("32")) + "\""}
	return validate(r, tps, nil)
}
