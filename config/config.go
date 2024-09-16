package config

import (
	_ "embed"
	"strings"

	"github.com/codevault-llc/humblebrag-api/parsers"
	regexp "github.com/wasilibs/go-re2"

	"github.com/spf13/viper"
)

type ViperConfig struct {
	Description string

	Rules []struct {
		ID          string
		Description string
		Regex       string
		Keywords    []string
	}

	Lists      []List
	ParserList []string
}

type Config struct {
	Rules      map[string]Rule
	Lists      map[string]List
	ParserList map[string]parsers.Parser
}

// Order the rules based on alphabetical order of the ID
func (vc *ViperConfig) OrderRules() []Rule {
	rules := make([]Rule, len(vc.Rules))

	for i, rule := range vc.Rules {
		rules[i] = Rule{
			Description: rule.Description,
			RuleID:      rule.ID,
			Regex:       regexp.MustCompile(rule.Regex),
			Keywords:    rule.Keywords,
		}
	}

	return rules
}
