package rules

import (
	"github.com/codevault-llc/humblebrag-api/cmd/generate/secrets"
	"github.com/codevault-llc/humblebrag-api/config"
)

func TrelloAccessToken() *config.Rule {
	// define rule
	r := config.Rule{
		RuleID:      "trello-access-token",
		Description: "Trello Access Token",
		Regex:       generateSemiGenericRegex([]string{"trello"}, `[a-zA-Z-0-9]{32}`, true),

		Keywords: []string{
			"trello",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("trello", secrets.NewSecret(`[a-zA-Z-0-9]{32}`)),
	}
	return validate(r, tps, nil)
}
