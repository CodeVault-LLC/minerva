package utils

import (
	"testing"

	"github.com/codevault-llc/humblebrag-api/pkg/types"
	regexp "github.com/wasilibs/go-re2"
)

func TestGenericScan(t *testing.T) {
	rule := types.Rule{
		Description: "Test rule",
		RuleID:      "test-rule",
		Keywords:    []string{"test"},
		Regex:       regexp.MustCompile(`(?i)test`), // Match "test" case-insensitively
	}

	content := "Hello, this is just a simple little fun test string which should work perfectly."

	script := types.FileRequest{
		Src:        "test",
		Content:    content,
		HashedBody: "test",
		FileSize:   0,
		FileType:   "text",
	}

	matches := GenericScan(rule, script)

	if len(matches) == 0 {
		t.Error("Expected some length to matches, got 0")
	}

	if matches[0].Match != "test" {
		t.Error("Expected match to be 'test', got", matches[0].Match)
	}
}

func TestFindMatchingLine(t *testing.T) {
	content := `Hello, this is just a simple little fun test string which should work perfectly.
This is the second line of the content.
This is the third line of the content.
This is the fourth line of the content.
`

	match := "test"

	line := findMatchingLine(content, match)

	if line != 1 {
		t.Error("Expected line 1, got", line)
	}

	match = "second"

	line = findMatchingLine(content, match)

	if line != 2 {
		t.Error("Expected line 2, got", line)
	}
}

func TestNewSecret(t *testing.T) {
	regex := "[a-z]{10}"
	secret := NewSecret(regex)

	if len(secret) != 10 {
		t.Error("Expected secret to be 10 characters long, got", len(secret))
	}
}
