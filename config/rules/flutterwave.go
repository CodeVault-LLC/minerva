package rules

import (
	regexp "github.com/wasilibs/go-re2"

	"github.com/codevault-llc/minerva/pkg/types"
	"github.com/codevault-llc/minerva/pkg/utils"
)

func FlutterwavePublicKey() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Detected a Finicity Public Key, potentially exposing public cryptographic operations and integrations.",
		RuleID:      "flutterwave-public-key",
		Regex:       regexp.MustCompile(`FLWPUBK_TEST-(?i)[a-h0-9]{32}-X`),
		Keywords:    []string{"FLWPUBK_TEST"},
	}

	// validate
	tps := []string{
		generateSampleSecret("flutterwavePubKey", "FLWPUBK_TEST-"+utils.NewSecret(hex("32"))+"-X"),
	}
	return validate(r, tps, nil)
}

func FlutterwaveSecretKey() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Identified a Flutterwave Secret Key, risking unauthorized financial transactions and data breaches.",
		RuleID:      "flutterwave-secret-key",
		Regex:       regexp.MustCompile(`FLWSECK_TEST-(?i)[a-h0-9]{32}-X`),
		Keywords:    []string{"FLWSECK_TEST"},
	}

	// validate
	tps := []string{
		generateSampleSecret("flutterwavePubKey", "FLWSECK_TEST-"+utils.NewSecret(hex("32"))+"-X"),
	}
	return validate(r, tps, nil)
}

func FlutterwaveEncKey() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Uncovered a Flutterwave Encryption Key, which may compromise payment processing and sensitive financial information.",
		RuleID:      "flutterwave-encryption-key",
		Regex:       regexp.MustCompile(`FLWSECK_TEST-(?i)[a-h0-9]{12}`),
		Keywords:    []string{"FLWSECK_TEST"},
	}

	// validate
	tps := []string{
		generateSampleSecret("flutterwavePubKey", "FLWSECK_TEST-"+utils.NewSecret(hex("12"))),
	}
	return validate(r, tps, nil)
}
