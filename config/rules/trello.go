package rules

import (
	"github.com/codevault-llc/minerva/pkg/types"
	"github.com/codevault-llc/minerva/pkg/utils"
)

func TrelloAccessToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "trello-access-token",
		Description: "Trello Access Token",
		Regex:       generateSemiGenericRegex([]string{"trello"}, `[a-zA-Z-0-9]{32}`, true),

		Keywords: []string{
			"trello",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("trello", utils.NewSecret(`[a-zA-Z-0-9]{32}`)),
	}
	return validate(r, tps, nil)
}
