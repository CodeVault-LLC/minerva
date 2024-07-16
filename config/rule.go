package config

import (
	"regexp"
)

type Rule struct {
	Description string

	RuleID string

	Regex *regexp.Regexp

	Keywords []string
}
