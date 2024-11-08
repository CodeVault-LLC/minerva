package rules

import (
	"fmt"

	"github.com/codevault-llc/minerva/pkg/types"
	"github.com/codevault-llc/minerva/pkg/utils"
)

func Authress() *types.Rule {
	// types.Rule Definition
	// (Note: When changes are made to this, rerun `go generate ./...` and commit the config/gitleaks.toml file
	r := types.Rule{
		Description: "Uncovered a possible Authress Service Client Access Key, which may compromise access control services and sensitive data.",
		RuleID:      "authress-service-client-access-key",
		Regex:       generateUniqueTokenRegex(`(?:sc|ext|scauth|authress)_[a-z0-9]{5,30}\.[a-z0-9]{4,6}\.acc[_-][a-z0-9-]{10,32}\.[a-z0-9+/_=-]{30,120}`, true),
		Keywords:    []string{"sc_", "ext_", "scauth_", "authress_"},
	}

	// validate
	// https://authress.io/knowledge-base/docs/authorization/service-clients/secrets-scanning/#1-detection
	service_client_id := "sc_" + alphaNumeric("10")
	access_key_id := alphaNumeric("4")
	account_id := "acc_" + alphaNumeric("10")
	signature_key := alphaNumericExtendedShort("40")

	tps := []string{
		generateSampleSecret("authress", utils.NewSecret(fmt.Sprintf(`%s\.%s\.%s\.%s`, service_client_id, access_key_id, account_id, signature_key))),
	}
	return validate(r, tps, nil)
}
