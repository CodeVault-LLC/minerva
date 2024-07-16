package config

import (
	_ "embed"
	"regexp"
	"strings"

	"github.com/spf13/viper"
)

//go:embed humblebrag.toml
var DefaultConfig string

type ViperConfig struct {
	Description string

	Rules []struct {
		ID string
		Description string
		Regex string
		Keywords []string
	}
}

type Config struct {
	Rules map[string]Rule
}

// Mainly read the config file and return the ViperConfig struct
func (vc *ViperConfig) ReadConfig() error {
	viper.SetConfigType("toml")
	viper.ReadConfig(strings.NewReader(DefaultConfig))

	err := viper.Unmarshal(vc)
	if err != nil {
		return err
	}

	return nil
}

// Order the rules based on alphabetical order of the ID
func (vc *ViperConfig) OrderRules() []Rule {
	rules := make([]Rule, len(vc.Rules))

	for i, rule := range vc.Rules {
		rules[i] = Rule{
			Description: rule.Description,
			RuleID: rule.ID,
			Regex: regexp.MustCompile(rule.Regex),
			Keywords: rule.Keywords,
		}
	}

	return rules
}
