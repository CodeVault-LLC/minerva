package utils

import (
	"log"
	"strings"

	"github.com/codevault-llc/humblebrag-api/config"
)

type Script struct {
	Src string `json:"src"`
	Content string `json:"content"`
}

func GenericScan(rule config.Rule, script Script) []Match {
	// Find all matches in the script content
	matches := rule.Regex.FindAllString(script.Content, -1)

	// Initialize the return structure
	var result RegexReturn
	result.Name = rule.RuleID
	result.Matches = make([]Match, 0, len(matches))

	// Process each match
	for _, match := range matches {
		if match != "" {
			log.Println("Match found: ", match)
			line := findMatchingLine(script.Content, match)
			result.Matches = append(result.Matches, Match{Match: match, Line: line, Source: script.Src})
		} else {
			log.Println("Match is empty")
			continue
		}
	}

	return result.Matches
}

// findMatchingLine returns the line containing the match in the content
func findMatchingLine(content, match string) int {
	lines := strings.Split(content, "\n")

	// Iterate over each line to find the one containing the match
	for _, line := range lines {
		if strings.Contains(line, match) {
			return strings.Index(content, line) + 1
		}
	}

	return 0
}
