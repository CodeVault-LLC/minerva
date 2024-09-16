package rules

import (
	regexp "github.com/wasilibs/go-re2"

	"github.com/codevault-llc/humblebrag-api/types"
	"github.com/codevault-llc/humblebrag-api/utils"
)

func GitHubPat() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Uncovered a GitHub Personal Access Token, potentially leading to unauthorized repository access and sensitive content exposure.",
		RuleID:      "github-pat",
		Regex:       regexp.MustCompile(`ghp_[0-9a-zA-Z]{36}`),
		Keywords:    []string{"ghp_"},
	}

	// validate
	tps := []string{
		generateSampleSecret("github", "ghp_"+utils.NewSecret(alphaNumeric("36"))),
	}
	return validate(r, tps, nil)
}

func GitHubFineGrainedPat() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Found a GitHub Fine-Grained Personal Access Token, risking unauthorized repository access and code manipulation.",
		RuleID:      "github-fine-grained-pat",
		Regex:       regexp.MustCompile(`github_pat_[0-9a-zA-Z_]{82}`),
		Keywords:    []string{"github_pat_"},
	}

	// validate
	tps := []string{
		generateSampleSecret("github", "github_pat_"+utils.NewSecret(alphaNumeric("82"))),
	}
	return validate(r, tps, nil)
}

func GitHubOauth() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Discovered a GitHub OAuth Access Token, posing a risk of compromised GitHub account integrations and data leaks.",
		RuleID:      "github-oauth",
		Regex:       regexp.MustCompile(`gho_[0-9a-zA-Z]{36}`),
		Keywords:    []string{"gho_"},
	}

	// validate
	tps := []string{
		generateSampleSecret("github", "gho_"+utils.NewSecret(alphaNumeric("36"))),
	}
	return validate(r, tps, nil)
}

func GitHubApp() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Identified a GitHub App Token, which may compromise GitHub application integrations and source code security.",
		RuleID:      "github-app-token",
		Regex:       regexp.MustCompile(`(ghu|ghs)_[0-9a-zA-Z]{36}`),
		Keywords:    []string{"ghu_", "ghs_"},
	}

	// validate
	tps := []string{
		generateSampleSecret("github", "ghu_"+utils.NewSecret(alphaNumeric("36"))),
		generateSampleSecret("github", "ghs_"+utils.NewSecret(alphaNumeric("36"))),
	}
	return validate(r, tps, nil)
}

func GitHubRefresh() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Detected a GitHub Refresh Token, which could allow prolonged unauthorized access to GitHub services.",
		RuleID:      "github-refresh-token",
		Regex:       regexp.MustCompile(`ghr_[0-9a-zA-Z]{36}`),
		Keywords:    []string{"ghr_"},
	}

	// validate
	tps := []string{
		generateSampleSecret("github", "ghr_"+utils.NewSecret(alphaNumeric("36"))),
	}
	return validate(r, tps, nil)
}
