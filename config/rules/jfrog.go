package rules

import (
	"fmt"

	"github.com/codevault-llc/minerva/pkg/types"
	"github.com/codevault-llc/minerva/pkg/utils"
)

func JFrogAPIKey() *types.Rule {
	keywords := []string{"jfrog", "artifactory", "bintray", "xray"}

	// Define types.Rule
	r := types.Rule{
		// Human readable description of the types.Rule
		Description: "Found a JFrog API Key, posing a risk of unauthorized access to software artifact repositories and build pipelines.",

		// Unique ID for the types.Rule
		RuleID: "jfrog-api-key",

		// Regex capture group for the actual secret

		// Regex used for detecting secrets. See regex section below for more details
		Regex: generateSemiGenericRegex(keywords, alphaNumeric("73"), true),

		// Keywords used for string matching on fragments (think of this as a prefilter)
		Keywords: keywords,
	}
	// validate
	tps := []string{
		fmt.Sprintf("--set imagePullSecretJfrog.password=%s", utils.NewSecret(alphaNumeric("73"))),
	}
	return validate(r, tps, nil)
}

func JFrogIdentityToken() *types.Rule {
	keywords := []string{"jfrog", "artifactory", "bintray", "xray"}

	// Define types.Rule
	r := types.Rule{
		// Human readable description of the types.Rule
		Description: "Discovered a JFrog Identity Token, potentially compromising access to JFrog services and sensitive software artifacts.",

		// Unique ID for the types.Rule
		RuleID: "jfrog-identity-token",

		// Regex capture group for the actual secret

		// Regex used for detecting secrets. See regex section below for more details
		Regex: generateSemiGenericRegex(keywords, alphaNumeric("64"), true),

		// Keywords used for string matching on fragments (think of this as a prefilter)
		Keywords: keywords,
	}

	// validate
	tps := []string{
		generateSampleSecret("jfrog", utils.NewSecret(alphaNumeric("64"))),
		generateSampleSecret("artifactory", utils.NewSecret(alphaNumeric("64"))),
		generateSampleSecret("bintray", utils.NewSecret(alphaNumeric("64"))),
		generateSampleSecret("xray", utils.NewSecret(alphaNumeric("64"))),
	}
	return validate(r, tps, nil)
}
