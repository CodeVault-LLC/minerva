package rules

import (
	regexp "github.com/wasilibs/go-re2"

	"github.com/codevault-llc/humblebrag-api/pkg/types"
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
)

func TeamsWebhook() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Uncovered a Microsoft Teams Webhook, which could lead to unauthorized access to team collaboration tools and data leaks.",
		RuleID:      "microsoft-teams-webhook",
		Regex: regexp.MustCompile(
			`https:\/\/[a-z0-9]+\.webhook\.office\.com\/webhookb2\/[a-z0-9]{8}-([a-z0-9]{4}-){3}[a-z0-9]{12}@[a-z0-9]{8}-([a-z0-9]{4}-){3}[a-z0-9]{12}\/IncomingWebhook\/[a-z0-9]{32}\/[a-z0-9]{8}-([a-z0-9]{4}-){3}[a-z0-9]{12}`),
		Keywords: []string{
			"webhook.office.com",
			"webhookb2",
			"IncomingWebhook",
		},
	}

	// validate
	tps := []string{
		"https://mycompany.webhook.office.com/webhookb2/" + utils.NewSecret(`[a-z0-9]{8}-([a-z0-9]{4}-){3}[a-z0-9]{12}@[a-z0-9]{8}-([a-z0-9]{4}-){3}[a-z0-9]{12}\/IncomingWebhook\/[a-z0-9]{32}\/[a-z0-9]{8}-([a-z0-9]{4}-){3}[a-z0-9]{12}`), // gitleaks:allow
	}
	return validate(r, tps, nil)
}
