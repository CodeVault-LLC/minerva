package rules

import (
	regexp "github.com/wasilibs/go-re2"

	"github.com/codevault-llc/humblebrag-api/types"
	"github.com/codevault-llc/humblebrag-api/utils"
)

func ShopifySharedSecret() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Found a Shopify shared secret, posing a risk to application authentication and e-commerce platform security.",
		RuleID:      "shopify-shared-secret",
		Regex:       regexp.MustCompile(`shpss_[a-fA-F0-9]{32}`),
		Keywords:    []string{"shpss_"},
	}

	// validate
	tps := []string{"shopifySecret := \"shpss_" + utils.NewSecret(hex("32")) + "\""}
	return validate(r, tps, nil)
}

func ShopifyAccessToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Uncovered a Shopify access token, which could lead to unauthorized e-commerce platform access and data breaches.",
		RuleID:      "shopify-access-token",
		Regex:       regexp.MustCompile(`shpat_[a-fA-F0-9]{32}`),
		Keywords:    []string{"shpat_"},
	}

	// validate
	tps := []string{"shopifyToken := \"shpat_" + utils.NewSecret(hex("32")) + "\""}
	return validate(r, tps, nil)
}

func ShopifyCustomAccessToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Detected a Shopify custom access token, potentially compromising custom app integrations and e-commerce data security.",
		RuleID:      "shopify-custom-access-token",
		Regex:       regexp.MustCompile(`shpca_[a-fA-F0-9]{32}`),
		Keywords:    []string{"shpca_"},
	}

	// validate
	tps := []string{"shopifyToken := \"shpca_" + utils.NewSecret(hex("32")) + "\""}
	return validate(r, tps, nil)
}

func ShopifyPrivateAppAccessToken() *types.Rule {
	// define types.Rule
	r := types.Rule{
		Description: "Identified a Shopify private app access token, risking unauthorized access to private app data and store operations.",
		RuleID:      "shopify-private-app-access-token",
		Regex:       regexp.MustCompile(`shppa_[a-fA-F0-9]{32}`),
		Keywords:    []string{"shppa_"},
	}

	// validate
	tps := []string{"shopifyToken := \"shppa_" + utils.NewSecret(hex("32")) + "\""}
	return validate(r, tps, nil)
}
