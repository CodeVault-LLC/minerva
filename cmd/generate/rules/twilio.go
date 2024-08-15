package rules

import (
	regexp "github.com/wasilibs/go-re2"

	"github.com/codevault-llc/humblebrag-api/cmd/generate/secrets"
	"github.com/codevault-llc/humblebrag-api/config"
)

func Twilio() *config.Rule {
	// define rule
	r := config.Rule{
		Description: "Found a Twilio API Key, posing a risk to communication services and sensitive customer interaction data.",
		RuleID:      "twilio-api-key",
		Regex:       regexp.MustCompile(`SK[0-9a-fA-F]{32}`),
		Keywords:    []string{"twilio"},
	}

	// validate
	tps := []string{
		"twilioAPIKey := \"SK" + secrets.NewSecret(hex("32")) + "\"",
	}
	return validate(r, tps, nil)
}
