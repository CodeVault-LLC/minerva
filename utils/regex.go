package utils

import (
	"log"
	"regexp"
	"strings"
)

type Script struct {
	Src string `json:"src"`
	Content string `json:"content"`
}

func GenericScan(pattern RegexPattern, script Script) []Match {
	// Compile the regex pattern
	re := regexp.MustCompile(pattern.Pattern)
	// Find all matches in the script content
	matches := re.FindAllString(script.Content, -1)

	// Initialize the return structure
	var result RegexReturn
	result.Name = pattern.Name
	result.Matches = make([]Match, 0, len(matches))

	// Process each match
	for _, match := range matches {
		if match != "" {
			log.Println("Match found: ", match)
			// Get the line containing the match
			line := findMatchingLine(script.Content, match)
			// Append the match details to the result
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
