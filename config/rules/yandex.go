package rules

import (
	"github.com/codevault-llc/minerva/pkg/types"
	"github.com/codevault-llc/minerva/pkg/utils"
)

func YandexAWSAccessToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "yandex-aws-access-token",
		Description: "Uncovered a Yandex AWS Access Token, potentially compromising cloud resource access and data security on Yandex Cloud.",
		Regex: generateSemiGenericRegex([]string{"yandex"},
			`YC[a-zA-Z0-9_\-]{38}`, true),
		Keywords: []string{
			"yandex",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("yandex",
			utils.NewSecret(`YC[a-zA-Z0-9_\-]{38}`)),
	}
	return validate(r, tps, nil)
}

func YandexAPIKey() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "yandex-api-key",
		Description: "Discovered a Yandex API Key, which could lead to unauthorized access to Yandex services and data manipulation.",
		Regex: generateSemiGenericRegex([]string{"yandex"},
			`AQVN[A-Za-z0-9_\-]{35,38}`, true),

		Keywords: []string{
			"yandex",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("yandex",
			utils.NewSecret(`AQVN[A-Za-z0-9_\-]{35,38}`)),
	}
	return validate(r, tps, nil)
}

func YandexAccessToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
		RuleID:      "yandex-access-token",
		Description: "Found a Yandex Access Token, posing a risk to Yandex service integrations and user data privacy.",
		Regex: generateSemiGenericRegex([]string{"yandex"},
			`t1\.[A-Z0-9a-z_-]+[=]{0,2}\.[A-Z0-9a-z_-]{86}[=]{0,2}`, true),

		Keywords: []string{
			"yandex",
		},
	}

	// validate
	tps := []string{
		generateSampleSecret("yandex",
			utils.NewSecret(`t1\.[A-Z0-9a-z_-]+[=]{0,2}\.[A-Z0-9a-z_-]{86}[=]{0,2}`)),
	}
	return validate(r, tps, nil)
}
