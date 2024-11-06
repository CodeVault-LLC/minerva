package rules

import (
	"fmt"
	"strings"

	"github.com/codevault-llc/minerva/pkg/types"
	"github.com/codevault-llc/minerva/pkg/utils"
	regexp "github.com/wasilibs/go-re2"

	"github.com/rs/zerolog/log"
)

const (
	// case insensitive prefix
	caseInsensitive = `(?i)`

	// identifier prefix (just an ignore group)
	identifierCaseInsensitivePrefix = `(?i:`
	identifierCaseInsensitiveSuffix = `)`
	identifierPrefix                = `(?:`
	identifierSuffix                = `)(?:[0-9a-z\-_\t .]{0,20})(?:[\s|']|[\s|"]){0,3}`

	// commonly used assignment operators or function call
	operator = `(?:=|>|:{1,3}=|\|\|:|<=|=>|:|\?=)`

	// boundaries for the secret
	// \x60 = `
	secretPrefixUnique = `\b(`
	secretPrefix       = `(?:'|\"|\s|=|\x60){0,5}(`
	secretSuffix       = `)(?:['|\"|\n|\r|\s|\x60|;]|$)`
)

func generateSemiGenericRegex(identifiers []string, secretRegex string, isCaseInsensitive bool) *regexp.Regexp {
	var sb strings.Builder
	// The identifiers should always be case-insensitive.
	// This is inelegant but prevents an extraneous `(?i:)` from being added to the pattern; it could be removed.
	if isCaseInsensitive {
		sb.WriteString(caseInsensitive)
		writeIdentifiers(&sb, identifiers)
	} else {
		sb.WriteString(identifierCaseInsensitivePrefix)
		writeIdentifiers(&sb, identifiers)
		sb.WriteString(identifierCaseInsensitiveSuffix)
	}
	sb.WriteString(operator)
	sb.WriteString(secretPrefix)
	sb.WriteString(secretRegex)
	sb.WriteString(secretSuffix)
	return regexp.MustCompile(sb.String())
}

func writeIdentifiers(sb *strings.Builder, identifiers []string) {
	sb.WriteString(identifierPrefix)
	sb.WriteString(strings.Join(identifiers, "|"))
	sb.WriteString(identifierSuffix)
}

func generateUniqueTokenRegex(secretRegex string, isCaseInsensitive bool) *regexp.Regexp {
	var sb strings.Builder
	if isCaseInsensitive {
		sb.WriteString(caseInsensitive)
	}
	sb.WriteString(secretPrefixUnique)
	sb.WriteString(secretRegex)
	sb.WriteString(secretSuffix)
	return regexp.MustCompile(sb.String())
}

func generateSampleSecret(identifier string, secret string) string {
	return fmt.Sprintf("%s_api_token = \"%s\"", identifier, secret)
}

func validate(r types.Rule, truePositives []string, falsePositives []string) *types.Rule {
	// normalize keywords like in the config package
	var keywords []string
	for _, k := range r.Keywords {
		keywords = append(keywords, strings.ToLower(k))
	}
	r.Keywords = keywords

	rules := make(map[string]types.Rule)
	rules[r.RuleID] = r
	for _, tp := range truePositives {
		if len(utils.GenericScan(r, types.FileRequest{
			Src:     "validate.js",
			Content: tp,
		})) != 1 {
			log.Fatal().Msgf("Failed to validate. For types.Rule ID [%s], true positive [%s] was not detected by regexp [%s]", r.RuleID, tp, r.Regex)
		}
	}
	for _, fp := range falsePositives {
		if len(utils.GenericScan(r, types.FileRequest{
			Src:     "validate.js",
			Content: fp,
		})) != 0 {
			log.Fatal().Msgf("Failed to validate. For types.Rule ID [%s], false positive [%s] was detected by regexp [%s]", r.RuleID, fp, r.Regex)
		}
	}
	return &r
}

func numeric(size string) string {
	return fmt.Sprintf(`[0-9]{%s}`, size)
}

func hex(size string) string {
	return fmt.Sprintf(`[a-f0-9]{%s}`, size)
}

func alphaNumeric(size string) string {
	return fmt.Sprintf(`[a-z0-9]{%s}`, size)
}

func alphaNumericExtendedShort(size string) string {
	return fmt.Sprintf(`[a-z0-9_-]{%s}`, size)
}

func alphaNumericExtended(size string) string {
	return fmt.Sprintf(`[a-z0-9=_\-]{%s}`, size)
}

func alphaNumericExtendedLong(size string) string {
	return fmt.Sprintf(`[a-z0-9\/=_\+\-]{%s}`, size)
}

func hex8_4_4_4_12() string {
	return `[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`
}
