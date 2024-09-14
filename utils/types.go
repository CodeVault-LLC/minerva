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
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Matches     []Match `json:"matches"`
}

type DiscordUser struct {
	Id            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Avatar        string `json:"avatar"`
	Verified      bool   `json:"verified"`
	Email         string `json:"email"`
	Flags         int    `json:"flags"`
}
