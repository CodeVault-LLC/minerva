package types

import regexp "github.com/wasilibs/go-re2"

type FingerprintType string

const (
	FingerprintTypeScript FingerprintType = "script"
)

type Fingerprint struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Regex       *regexp.Regexp  `json:"regex"`
	Type        FingerprintType `json:"type"`
	Keywords    []string        `json:"keywords"`
}
