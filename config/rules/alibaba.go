package rules

import (
	"github.com/codevault-llc/humblebrag-api/types"
	"github.com/codevault-llc/humblebrag-api/utils"
)

func AlibabaAccessKey() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Detected an Alibaba Cloud AccessKey ID, posing a risk of unauthorized cloud resource access and potential data compromise.",
		RuleID:      "alibaba-access-key-id",
		Regex:       generateUniqueTokenRegex(`(LTAI)(?i)[a-z0-9]{20}`, true),
		Keywords:    []string{"LTAI"},
	}

	// validate
	tps := []string{
		"alibabaKey := \"LTAI" + utils.NewSecret(hex("20")) + "\"",
	}
	return validate(r, tps, nil)
}

// TODO
func AlibabaSecretKey() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Discovered a potential Alibaba Cloud Secret Key, potentially allowing unauthorized operations and data access within Alibaba Cloud.",
		RuleID:      "alibaba-secret-key",
		Regex: generateSemiGenericRegex([]string{"alibaba"},
			alphaNumeric("30"), true),

		Keywords: []string{"alibaba"},
	}

	// validate
	tps := []string{
		generateSampleSecret("alibaba", utils.NewSecret(alphaNumeric("30"))),
	}
	return validate(r, tps, nil)
}
