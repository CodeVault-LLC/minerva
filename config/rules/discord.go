package rules

import (
	"github.com/codevault-llc/humblebrag-api/pkg/types"
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
)

func DiscordAPIToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Detected a Discord API key, potentially compromising communication channels and user data privacy on Discord.",
		RuleID:      "discord-api-token",
		Regex:       generateSemiGenericRegex([]string{"discord"}, hex("64"), true),
		Keywords:    []string{"discord"},
	}

	// validate
	tps := []string{
		generateSampleSecret("discord", utils.NewSecret(hex("64"))),
	}
	return validate(r, tps, nil)
}

func DiscordClientID() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Identified a Discord client ID, which may lead to unauthorized integrations and data exposure in Discord applications.",
		RuleID:      "discord-client-id",
		Regex:       generateSemiGenericRegex([]string{"discord"}, numeric("18"), true),
		Keywords:    []string{"discord"},
	}

	// validate
	tps := []string{
		generateSampleSecret("discord", utils.NewSecret(numeric("18"))),
	}
	return validate(r, tps, nil)
}

func DiscordClientSecret() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Discovered a potential Discord client secret, risking compromised Discord bot integrations and data leaks.",
		RuleID:      "discord-client-secret",
		Regex:       generateSemiGenericRegex([]string{"discord"}, alphaNumericExtended("32"), true),
		Keywords:    []string{"discord"},
	}

	// validate
	tps := []string{
		generateSampleSecret("discord", utils.NewSecret(numeric("32"))),
	}
	return validate(r, tps, nil)
}
