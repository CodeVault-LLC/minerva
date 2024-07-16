package utils

type RegexPattern struct {
	Pattern string `json:"pattern"`
	Name    string
}

type Match struct {
	Match  string `json:"match"`
	Line   int    `json:"line"`
	Source string `json:"source"`
}

type RegexReturn struct {
	Name    string  `json:"name"`
	Matches []Match `json:"matches"`
}
