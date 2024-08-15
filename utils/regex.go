package utils

import (
	"log"
	"strings"

	regexp "github.com/wasilibs/go-re2"

	"github.com/codevault-llc/humblebrag-api/config"
)

type Script struct {
	Src     string `json:"src"`
	Content string `json:"content"`
}

func GenericScan(rule config.Rule, script Script) []Match {
	re, err := regexp.Compile(rule.Regex.String())
	if err != nil {
		log.Fatalf("Failed to compile regex: %v", err)
	}

	var result RegexReturn
	result.Name = rule.RuleID
	result.Description = rule.Description
	result.Matches = make([]Match, 0)

	matches := re.FindAllIndex([]byte(script.Content), -1)

	for _, match := range matches {
		matchStr := script.Content[match[0]:match[1]]

		if matchStr != "" {
			line := findMatchingLine(script.Content, matchStr)
			result.Matches = append(result.Matches, Match{Match: matchStr, Line: line, Source: script.Src})
		}
	}

	return result.Matches
}

// findMatchingLine returns the line containing the match in the content
func findMatchingLine(content, match string) int {
	lines := strings.Split(content, "\n")

	for i, line := range lines {
		if strings.Contains(line, match) {
			return i + 1
		}
	}

	return 0
}
