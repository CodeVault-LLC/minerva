package rules

import (
	"github.com/codevault-llc/humblebrag-api/pkg/types"
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
)

func OpenAI() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "openai-api-key",
		Description: "Found an OpenAI API Key, posing a risk of unauthorized access to AI services and data manipulation.",
		Regex:       generateUniqueTokenRegex(`sk-[a-zA-Z0-9]{20}T3BlbkFJ[a-zA-Z0-9]{20}`, true),

		Keywords: []string{
			"T3BlbkFJ",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("openaiApiKey", "sk-"+utils.NewSecret(alphaNumeric("20"))+"T3BlbkFJ"+utils.NewSecret(alphaNumeric("20"))),
	}
	return validate(r, tps, nil)
}
