package rules

import (
	"github.com/codevault-llc/humblebrag-api/cmd/generate/secrets"
	"github.com/codevault-llc/humblebrag-api/config"
)

func OpenAI() *config.Rule {
	// define rule
	r := config.Rule{
		RuleID:      "openai-api-key",
		Description: "Found an OpenAI API Key, posing a risk of unauthorized access to AI services and data manipulation.",
		Regex:       generateUniqueTokenRegex(`sk-[a-zA-Z0-9]{20}T3BlbkFJ[a-zA-Z0-9]{20}`, true),

		Keywords: []string{
			"T3BlbkFJ",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("openaiApiKey", "sk-"+secrets.NewSecret(alphaNumeric("20"))+"T3BlbkFJ"+secrets.NewSecret(alphaNumeric("20"))),
	}
	return validate(r, tps, nil)
}
