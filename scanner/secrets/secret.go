package secrets

import (
	"github.com/codevault-llc/humblebrag-api/models"
	"github.com/codevault-llc/humblebrag-api/utils"
)

// make a list of regexes to match
var regexes = []utils.RegexPattern{
	{
		Name: "Sensitive Data Keywords",
		Pattern: `secret\s*:\s*['"][^'"]+`,
	},
	{
		Name: "Google API Key",
		Pattern: `AIza[0-9A-Za-z\-_]{35}`,
	},
	{
		Name: "AWS Access Key ID",
		Pattern: `AKIA[0-9A-Z]{16}`,
	},
	{
		Name: "Google Analytics Tracking ID",
		Pattern: `UA-[0-9]+-[0-9]+`,
	},
	{
		Name: "URL Matching",
		Pattern: `https?:\/\/[^\s/$.?#].[^\s"]*`,
	},
	{
		Name: "HTTP Matching",
		Pattern: `http?:\/\/[^\s/$.?#].[^\s"]*`,
	},
	{
		Name: "HTTPS Matching",
		Pattern: `https?:\/\/[^\s/$.?#].[^\s"]*`,
	},
	{
		Name: "SSH URL Matching",
		Pattern: `ssh:\/\/[^\s/$.?#].[^\s"]*`,
	},
	{
		Name: "FTP URL Matching",
		Pattern: `ftp:\/\/[^\s/$.?#].[^\s"]*`,
	},
	{
		Name: "Slack Webhook URL",
		Pattern: `https:\/\/hooks.slack.com\/services\/T[A-Z0-9]{10}\/B[A-Z0-9]{10}\/[A-Za-z0-9]{24}`,
	},
	{
		Name: "Slack Token",
		Pattern: `xox[baprs]-[0-9]{12}-[0-9]{12}-[0-9a-zA-Z]{24}`,
	},
}

func ScanSecrets(script []models.ScriptRequest) []utils.RegexReturn {
	var results []utils.RegexReturn

	// Loop through each secret and script, then once finished for the script append the results to the results slice
	for _, pattern := range regexes {
		for _, script := range script {
			matches := utils.GenericScan(pattern, utils.Script(script))
			if len(matches) > 0 {
				result := utils.RegexReturn{
					Name:    pattern.Name,
					Matches: matches,
				}
				results = append(results, result)
			}
		}
	}

	return results
}
