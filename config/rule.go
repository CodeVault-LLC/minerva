package config

import regexp "github.com/wasilibs/go-re2"

type Rule struct {
	Description string

	RuleID string

	Regex *regexp.Regexp

	Keywords []string
}
